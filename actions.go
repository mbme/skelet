package main

import (
	"encoding/json"
	"errors"

	"log"

	s "github.com/mbme/skelet/storage"
)

//ActionType type of action
type ActionType string

//ActionParams action raw parameters
type ActionParams json.RawMessage

//Possible actions
const (
	AtomsListReq ActionType = "req-atoms-list"
	AtomsList               = "atoms-list"
	AtomReq                 = "req-atom"
	Atom                    = "atom"
	NoType                  = ""
)

var (
	ErrorNoHandler = errors.New("no handler for action")
	ErrorNotFound  = errors.New("atom not found")
	ErrorBadParams = errors.New("malformed action params")
)

var storage = &s.VirtualStorage{}

type actionHandler func(*ActionParams) (ActionType, any, error)

type atomInfo struct {
	ID   s.AtomId   `json:"id"`
	Type s.AtomType `json:"type"`
	Name string     `json:"name"`
}

func toAtomInfo(atom *s.Atom) *atomInfo {
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
			infos[i] = toAtomInfo(atom)
		}

		return AtomsList, infos, nil
	},

	AtomReq: func(params *ActionParams) (ActionType, any, error) {
		link := &s.AtomLink{}
		if err := params.readAs(link); err != nil {
			log.Printf("error parsing params: %v", err)
			return NoType, nil, ErrorBadParams
		}

		if !link.IsValid() {
			log.Println("error parsing params: not valid atom link")
			return NoType, nil, ErrorBadParams
		}

		atom := storage.GetAtom(link)
		if atom == nil {
			log.Printf("atom not found: %v", link)
			return NoType, nil, ErrorNotFound
		}

		return Atom, atom, nil
	},
}

// HandleAction handle client action and produce own action
func HandleAction(actionType ActionType, params *ActionParams) (ActionType, any, error) {
	handler, ok := handlers[actionType]

	if !ok {
		return NoType, nil, ErrorNoHandler
	}

	return handler(params)
}
