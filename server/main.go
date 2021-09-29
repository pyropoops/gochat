package main

import (
	"embed"
	"fmt"
	"gochat/server/chat"
	"gochat/server/panel"
	"io/fs"
	"log"
	"net/http"
	"net"
	"strconv"

	"github.com/gorilla/mux"
)

const PORT = 80

//go:embed views
var embeds embed.FS

type App struct {
	ChatServer chat.Server
	Router     *mux.Router
}

func main() {
	app := App{
		ChatServer: chat.NewServer(),
		Router:     mux.NewRouter().StrictSlash(false),
	}

	views, err := fs.Sub(embeds, "views")
	if err != nil {
		panic(err)
	}

	app.Router.HandleFunc("/ws", app.ChatServer.HandleConnection)
	app.Router.PathPrefix("/").Handler(http.FileServer(http.FS(views)))

	logOutbound()

	controlPanel := panel.NewPanel(app.ChatServer)
	go controlPanel.Start()

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), app.Router); err != nil {
		log.Fatal(err)
	}
}

func logOutbound() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP
	fmt.Printf("Listening on: %s:%d\n", ip.String(), PORT)
}
