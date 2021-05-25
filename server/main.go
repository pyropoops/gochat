package main

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"gochat/server/chat"
	"gochat/server/panel"
	"io/fs"
	"log"
	"net/http"
	"strconv"
)

const PORT = 8000

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

	fmt.Printf("Listening on port %d...\n", PORT)

	controlPanel := panel.NewPanel(app.ChatServer)
	go controlPanel.Start()

	if err := http.ListenAndServe(":"+strconv.Itoa(PORT), app.Router); err != nil {
		log.Fatal(err)
	}
}