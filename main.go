package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// echo
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			log.Println(err)
			return
		}
		if err := w.Close(); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "go-skelet!"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8081,
			Usage: "websockets port",
		},
	}

	app.Action = func(c *cli.Context) {
		var port = c.String("port")
		fmt.Printf("listening on port %v\n", port)

		http.HandleFunc("/ws", handler)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
