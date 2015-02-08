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
	AtomUpdate              = "atom-update"
	NoType                  = ""
)

var (
	ErrorNoHandler = errors.New("no handler for action")
	ErrorBadParams = errors.New("malformed action params")
)

var storage = s.NewStorage()

type actionHandler func(*ActionParams) (ActionType, any, error)

type atomInfo struct {
	ID   s.AtomID   `json:"id"`
	Type s.AtomType `json:"type"`
	Name string     `json:"name"`
}

func toAtomInfo(atom *s.Atom) *atomInfo {
	return &atomInfo{
		ID:   *atom.ID,
		Type: *atom.Type,
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
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return NoType, nil, ErrorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return NoType, nil, ErrorBadParams
		}

		atom, err := storage.GetAtom(id)
		if err == nil {
			return NoType, nil, err
		}

		return Atom, atom, nil
	},

	AtomUpdate: func(params *ActionParams) (ActionType, any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return NoType, nil, ErrorBadParams
		}

		if !atom.IsValid() {
			log.Println("error parsing params: can't parse atom")
			return NoType, nil, ErrorBadParams
		}

		err := storage.UpdateAtom(atom)

		return NoType, nil, err
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
