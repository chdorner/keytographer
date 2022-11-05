package live

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/chdorner/keymap-render/internal/renderer"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var (
	//go:embed live.tpl
	tplLiveSrc string
)

type Server struct {
	debug bool
	host  string
	port  int
	r     renderer.Renderer

	mux     *http.ServeMux
	tplLive *template.Template

	clients     map[*websocket.Conn]bool
	broadcaster chan []byte
	upgrader    websocket.Upgrader

	watchFile string
}

func NewServer(r renderer.Renderer, watchFile, host string, port int, debug bool) (*Server, error) {
	tplLive, err := template.New("live").Parse(tplLiveSrc)
	if err != nil {
		return nil, err
	}

	s := &Server{
		debug: debug,
		host:  host,
		port:  port,
		r:     r,

		mux:     http.NewServeMux(),
		tplLive: tplLive,

		clients:     make(map[*websocket.Conn]bool),
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
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleWebsocketConnections(w http.ResponseWriter, req *http.Request) {
	ws, err := s.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	s.clients[ws] = true

	s.broadcaster <- s.r.Render()

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		log.Println(string(p))
	}
}

func (s *Server) watch(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				fmt.Println("hello")
				return
			}
			if event.Has(fsnotify.Write) {
				s.broadcaster <- s.r.Render()
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
