package main

import (
	"errors"

	s "github.com/mbme/skelet/storage"
)

var storage = &s.VirtualStorage{}

type actionHandler func(*ActionParams) (ActionType, any, error)

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
	AtomsListReq: func(_ *ActionParams) (ActionType, any, error) {
		atoms := storage.GetAtoms()
		infos := make([]*atomInfo, len(atoms))
		for i, atom := range atoms {
			infos[i] = newAtomInfo(atom)
		}

		return AtomsList, infos, nil
	},
}

var (
	ErrorNoHandler = errors.New("no handler for action")
)

// HandleAction handle client action and produce own action
func HandleAction(actionType ActionType, params *ActionParams) (ActionType, any, error) {
	handler, ok := handlers[actionType]

	if !ok {
		return NoType, nil, ErrorNoHandler
	}

	return handler(params)
}
