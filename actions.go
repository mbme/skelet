package main

import (
	"encoding/json"

	"fmt"

	s "github.com/mbme/skelet/storage"
)

var storage = &s.VirtualStorage{}

func toActionParams(data interface{}) (ActionParams, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := json.RawMessage(res)

	return &msg, nil
}

type actionHandler func(ActionParams) (ActionType, ActionParams, error)

type atomInfo struct {
	ID   s.AtomId   `json:"id"`
	Type s.AtomType `json:"type"`
	Name string     `json:"name"`
}

func newAtomInfo(atom *s.Atom) *atomInfo {
	return &atomInfo{
		ID:   atom.ID,
		Type: atom.Type,
		Name: atom.Name,
	}
}

var handlers = map[ActionType]actionHandler{
	RecordsListReq: func(_ ActionParams) (ActionType, ActionParams, error) {
		atoms := storage.GetAtoms()
		infos := make([]*atomInfo, len(atoms))
		for i, atom := range atoms {
			infos[i] = newAtomInfo(atom)
		}

		raw, err := toActionParams(infos)

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
