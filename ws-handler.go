package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type any interface{}

func (p *ActionParams) readAs(v any) error {
	return json.Unmarshal(*p, v)
}

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

func badRequestResp(errType string) *ActionResultWrapper {
	return &ActionResultWrapper{BadRequest, errType}
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
			writeResponse(conn, badRequestResp("can't parse"))
			continue
		}

		if req.Type == NoType {
			log.Printf("no type in request %v", req)
			writeResponse(conn, badRequestResp("missing type"))
			continue
		}

		params := ActionParams(req.Params)
		respType, respParams, err := HandleAction(req.Type, &params)

		if err != nil {
			log.Printf("%v -> %v", req.Type, err)
			writeResponse(conn, badRequestResp(err.Error()))
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

		writeResponse(conn, resp)
	}
}

func writeResponse(conn *websocket.Conn, resp *ActionResultWrapper) {
	if err := conn.WriteJSON(resp); err != nil {
		log.Printf("can't write response: %v", err)
	}
}
