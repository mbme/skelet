package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type any interface{}

func (p *RequestParams) readAs(v any) error {
	return json.Unmarshal(*p, v)
}

// JSON RPC
type requestID uint32

type request struct {
	Method RequestMethod   `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     *requestID      `json:"id"`
}

type response struct {
	Result any        `json:"result"`
	Error  error      `json:"error"`
	ID     *requestID `json:"id"`
}

// HTTP websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// allow all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//WsHandler websocket connection handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("connection: open")

	for {
		// parse request
		req := &request{}
		if err = conn.ReadJSON(req); err != nil {
			if err == io.EOF {
				log.Println("connection: closed")
				return
			}

			log.Printf("can't parse message: %v", err)
			return // close websocket
		}

		if req.Method == NoType || req.ID == nil {
			log.Printf("bad request %v", req)
			return // close websocket
		}

		resp := &response{
			ID: req.ID,
		}

		// process request
		params := RequestParams(req.Params)
		resp.Result, resp.Error = ProcessRequest(req.Method, &params)

		if resp.Error == nil {
			log.Printf("$v -> ok")
		} else {
			log.Printf("%v -> error: %v", req.Method, resp.Error)
		}

		// write response
		if err = conn.WriteJSON(resp); err != nil {
			log.Printf("can't write response: %v", err)
		}
	}
}
