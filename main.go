package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

	"io/ioutil"

	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/gorilla/websocket"
)

//ActionType type of action
type ActionType string

//Possible action types
const (
	RecordsListReq ActionType = "req-records-list"
	RecordsList               = "records-list"
	NoType                    = ""
)

//Action action
type Action struct {
	Type   ActionType       `json:"action"`
	Params *json.RawMessage `json:"params"`
}

//NoAction means do nothing
var NoAction = &Action{
	Type:   NoType,
	Params: nil,
}

func getRecordsList() *Action {
	return &Action{
		Type:   RecordsList,
		Params: nil,
	}
}

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
	for {
		// open reader
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		// skip binary messages
		if messageType != websocket.TextMessage {
			log.Fatalln("websocket: received binary message")
			return
		}

		// read request
		data, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println(err)
			return
		}

		// parse request
		var req *Action
		err = json.Unmarshal(data, &req)
		if err != nil {
			log.Println(err)
			return
		}

		// handle request
		resp := NoAction
		switch req.Type {
		case RecordsListReq:
			resp = getRecordsList()
		case NoType:
			log.Fatalf("no type in request %s\n", string(data))
		}

		// write response if required
		if resp != NoAction {
			log.Printf("%v -> %v\n", req.Type, resp.Type)

			// open writer
			w, err := conn.NextWriter(messageType)
			if err != nil {
				log.Println(err)
				return
			}

			// serialize data
			data, err := json.Marshal(resp)
			if err != nil {
				log.Println(err)
				return
			}

			// write response
			if _, err := w.Write(data); err != nil {
				log.Println(err)
				return
			}

			// close writer
			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
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
