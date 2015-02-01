package main

import (
	"encoding/json"

	"fmt"

	"github.com/mbme/skelet/storage"
)

var s = &storage.VirtualStorage{}

func toActionParams(data interface{}) (ActionParams, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := json.RawMessage(res)

	return &msg, nil
}

type actionHandler func(ActionParams) (ActionType, ActionParams, error)

var handlers = map[ActionType]actionHandler{
	RecordsListReq: func(_ ActionParams) (ActionType, ActionParams, error) {
		raw, err := toActionParams(s.GetAtoms())

		if err != nil {
			return NoType, nil, err
		}

		return RecordsList, raw, nil
	},
}

// HandleAction handle client action and produce own action
func HandleAction(actionType ActionType, params ActionParams) (ActionType, ActionParams, error) {
	handler, ok := handlers[actionType]

	if !ok {
		return NoType, nil, fmt.Errorf("no handler for action %s", actionType)
	}

	return handler(params)
}
