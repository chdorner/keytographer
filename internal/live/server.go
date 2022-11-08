package live

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	//go:embed live.tpl
	tplLiveSrc string
)

type Server struct {
	ctx      context.Context
	debug    bool
	host     string
	port     int
	renderer keytographer.Renderer

	mux     *http.ServeMux
	tplLive *template.Template

	clients     map[*websocket.Conn]bool
	broadcaster chan *RenderMessage
	upgrader    websocket.Upgrader

	watchFile string
}

type RenderMessage struct {
	Name     string         `json:"name"`
	Keyboard string         `json:"keyboard"`
	Layers   []MessageLayer `json:"layers"`
}

type MessageLayer struct {
	Name string `json:"name"`
	Svg  string `json:"svg"`
}

func NewServer(ctx context.Context, renderer keytographer.Renderer, watchFile, host string, port int) (*Server, error) {
	tplLive, err := template.New("live").Parse(tplLiveSrc)
	if err != nil {
		return nil, err
	}

	s := &Server{
		ctx:      ctx,
		debug:    ctx.Value("debug").(bool),
		host:     host,
		port:     port,
		renderer: renderer,

		mux:     http.NewServeMux(),
		tplLive: tplLive,

		clients: make(map[*websocket.Conn]bool),
		// TODO: also send config parsing / rendering errors back to browser
		broadcaster: make(chan *RenderMessage),
		upgrader:    websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},

		watchFile: watchFile,
	}

	s.mux.HandleFunc("/", s.liveHandler)
	s.mux.HandleFunc("/ws", s.handleWebsocketConnections)

	return s, nil
}

func (s *Server) ListenAndServe() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(s.watchFile)
	if err != nil {
		return err
	}

	go s.watch(watcher)
	go s.handlePushes()

	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s.mux)
}

func (s *Server) liveHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"debug":        s.debug,
		"websocketURL": fmt.Sprintf("ws://%s:%d/ws", s.host, s.port),
	}
	if err := s.tplLive.Execute(w, data); err != nil {
		logrus.WithField("err", err).Warn("failed to render live html template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleWebsocketConnections(w http.ResponseWriter, req *http.Request) {
	ws, err := s.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.WithField("err", err).Warn("failed to upgrade websocket connection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	s.clients[ws] = true

	output, _ := s.render()
	s.broadcaster <- output

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (s *Server) render() (*RenderMessage, error) {
	data, err := keytographer.Load(s.watchFile)
	if err != nil {
		return nil, err
	}

	err = keytographer.Validate(data)
	if err != nil {
		return nil, err
	}

	config, err := keytographer.Parse(data)
	if err != nil {
		return nil, err
	}

	msg := &RenderMessage{
		Name:     config.Name,
		Keyboard: config.Keyboard,
	}

	layers, err := s.renderer.RenderAllLayers(config)
	if err != nil {
		return nil, err
	}
	for _, layer := range layers {
		msg.Layers = append(msg.Layers, MessageLayer{
			Name: layer.Name,
			Svg:  string(layer.Svg),
		})
	}

	return msg, nil
}

func (s *Server) watch(watcher *fsnotify.Watcher) {
	//nolint:gosimple
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				output, _ := s.render()
				s.broadcaster <- output
			}
		}
	}
}

func (s *Server) handlePushes() {
	for {
		msg := <-s.broadcaster

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		for client := range s.clients {
			err = client.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}
