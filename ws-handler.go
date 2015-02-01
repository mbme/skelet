package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//ActionType type of action
type ActionType string

//ActionParams action raw parameters
type ActionParams *json.RawMessage

//Possible actions
const (
	AtomsListReq ActionType = "req-atoms-list"
	AtomsList               = "atoms-list"
	NoType                  = ""
)

//ActionWrapper action
type ActionWrapper struct {
	Type   ActionType   `json:"action"`
	Params ActionParams `json:"params"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//WsHandler websocket connection handler
func WsHandler(w http.ResponseWriter, r *http.Request) {
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
		var req *ActionWrapper
		err = json.Unmarshal(data, &req)
		if err != nil {
			log.Println(err)
			return
		}

		if req.Type == NoType {
			log.Fatalf("no type in request %s\n", string(data))
			continue
		}

		respType, respParams, err := HandleAction(req.Type, req.Params)

		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%v -> %v\n", req.Type, respType)

		if respType == NoType {
			continue
		}

		// write response

		// open writer
		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Println(err)
			return
		}

		resp := &ActionWrapper{
			Type:   respType,
			Params: respParams,
		}

		// serialize data
		data, err = json.Marshal(resp)
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
