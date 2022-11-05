package main

import (
	"flag"
	"log"
	"os"

	"github.com/chdorner/keymap-render/internal/live"
	"github.com/chdorner/keymap-render/internal/renderer"
)

func main() {
	var debug bool
	var configFile string
	var host string
	var port int

	flag.BoolVar(&debug, "d", false, "Enable debug mode.")
	flag.StringVar(&host, "h", "localhost", "Specify on which host to run the live server on.")
	flag.IntVar(&port, "p", 8080, "Specify on which port to run the live server on.")
	flag.StringVar(&configFile, "w", "", "Specify path to keymap file to watch for changes.")

	flag.Parse()

	if configFile == "" {
		log.Fatal("Need to specify path to keymap file with -w")
	}
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Specified file does not exist.")
	}

	if debug {
		log.Println("debug mode turned on.")
	}

	renderer := renderer.NewRenderer(&renderer.RenderConfig{})

	log.Printf("starting server on %s:%d\n", host, port)
	server, err := live.NewServer(renderer, configFile, host, port, debug)
	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
