package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/windler/godepg/config"
)

type webServer struct {
	root string
	port int64
}

func StartWebServer(port int64) {
	home := os.Getenv("GODEPG_HOME")
	if home == "" {
		home = config.GetDefaultHomeDir()
	}

	ws := &webServer{
		port: port,
		root: home,
	}

	ws.serve()
}

func (ws *webServer) serve() {
	fs := http.FileServer(http.Dir(ws.root))
	http.Handle("/", fs)

	fmt.Println("Started webserver on port " + strconv.FormatInt(ws.port, 10) + "...")
	fmt.Println("http://localhost:" + strconv.FormatInt(ws.port, 10))

	err := http.ListenAndServe(":"+strconv.FormatInt(ws.port, 10), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
