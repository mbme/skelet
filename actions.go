package main

import (
	"encoding/json"

	"fmt"

	"github.com/mbme/skelet/storage"
)

//ActionType type of action
type ActionType string

//Possible action types
const (
	RecordsListReq ActionType = "req-records-list"
	RecordsList               = "records-list"
	NoType                    = ""
)

var s = &storage.VirtualStorage{}

func toRawMessage(data interface{}) (*json.RawMessage, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := json.RawMessage(res)

	return &msg, err
}

type actionHandler func(*json.RawMessage) (ActionType, *json.RawMessage, error)

var handlers = map[ActionType]actionHandler{
	RecordsListReq: func(_ *json.RawMessage) (ActionType, *json.RawMessage, error) {
		raw, err := toRawMessage(s.GetAtoms())

		if err != nil {
			return NoType, nil, err
		}

		return RecordsList, raw, nil
	},
}

// HandleAction handle client action and produce own action
func HandleAction(actionType ActionType, params *json.RawMessage) (ActionType, *json.RawMessage, error) {
	handler, ok := handlers[actionType]

	if !ok {
		return NoType, nil, fmt.Errorf("no handler for action %s", actionType)
	}

	return handler(params)
}
