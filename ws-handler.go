package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type any interface{}

//ActionType type of action
type ActionType string

//ActionParams action raw parameters
type ActionParams json.RawMessage

func (p *ActionParams) ReadAs(v any) error {
	return json.Unmarshal(*p, v)
}

//Possible actions
const (
	AtomsListReq ActionType = "req-atoms-list"
	AtomsList               = "atoms-list"
	AtomReq                 = "req-atom"
	Atom                    = "atom"
	NoType                  = ""
)

//ActionWrapper action
type ActionWrapper struct {
	Type   ActionType      `json:"action"`
	Params json.RawMessage `json:"params"`
}

//ActionResultWrapper action result
type ActionResultWrapper struct {
	Type   ActionType `json:"action"`
	Params any        `json:"params"`
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
	log.Println("connection: open")

	for {
		// parse request
		req := &ActionWrapper{}
		if err = conn.ReadJSON(req); err != nil {
			if err == io.EOF {
				log.Println("connection: closed")
				return
			}

			log.Printf("can't parse message: %v", err)
			return
		}

		if req.Type == NoType {
			log.Printf("no type in request %v", req)
			continue
		}

		params := ActionParams(req.Params)
		respType, respParams, err := HandleAction(req.Type, &params)

		if err != nil {
			log.Printf("%v -> %v", req.Type, err)
			continue
		}

		log.Printf("%v -> %v", req.Type, respType)

		if respType == NoType {
			continue
		}

		// write response
		resp := &ActionResultWrapper{
			Type:   respType,
			Params: respParams,
		}

		// write response
		if err = conn.WriteJSON(resp); err != nil {
			log.Printf("can't write response: %v", err)
			continue
		}
	}
}
