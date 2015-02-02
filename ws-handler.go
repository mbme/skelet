package main

import (
	"encoding/json"
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
	for {
		// parse request
		req := &ActionWrapper{}
		if err = conn.ReadJSON(req); err != nil {
			log.Printf("can't parse message: %v\n", err)
			continue
		}

		if req.Type == NoType {
			log.Printf("no type in request %v\n", req)
			continue
		}

		params := ActionParams(req.Params)
		respType, respParams, err := HandleAction(req.Type, &params)

		if err != nil {
			log.Printf("%v -> %v\n", req.Type, err)
			continue
		}

		log.Printf("%v -> %v\n", req.Type, respType)

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
			log.Printf("can't write response: %v\n", err)
			continue
		}
	}
}
