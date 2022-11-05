package live

import (
	"context"
	_ "embed"
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
	broadcaster chan []byte
	upgrader    websocket.Upgrader

	watchFile string
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
		broadcaster: make(chan []byte),
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
	watcher.Add(s.watchFile)
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

func (s *Server) render() ([]byte, error) {
	config, err := keytographer.Parse(s.watchFile)
	if err != nil {
		return nil, err
	}

	return s.renderer.Render(config), nil
}

func (s *Server) watch(watcher *fsnotify.Watcher) {
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
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}
