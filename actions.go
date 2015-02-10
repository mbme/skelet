package main

import (
	"encoding/json"
	"errors"

	"log"

	"strings"

	s "github.com/mbme/skelet/storage"
)

//ActionType type of action
type ActionType string

//ActionParams action raw parameters
type ActionParams json.RawMessage

//Possible actions
const (
	AtomsListRead ActionType = "atoms-list-read"

	AtomRead   = "atom-read"
	AtomCreate = "atom-create"
	AtomUpdate = "atom-update"
	AtomDelete = "atom-delete"

	AtomsList = "atoms-list"
	Atom      = "atom"

	NoType = ""
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

func getAtomsList() []*atomInfo {
	atoms := storage.GetAtoms()
	infos := make([]*atomInfo, len(atoms))
	for i, atom := range atoms {
		infos[i] = toAtomInfo(atom)
	}

	return infos
}

var handlers = map[ActionType]actionHandler{
	AtomsListRead: func(_ *ActionParams) (ActionType, any, error) {
		return AtomsList, getAtomsList(), nil
	},

	AtomRead: func(params *ActionParams) (ActionType, any, error) {
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
		if err != nil {
			log.Printf("can't find atom %s", id)
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

		if atom.ID == nil || atom.Type == nil || !atom.Type.IsValid() {
			log.Println("error parsing params: bad atom")
			return NoType, nil, ErrorBadParams
		}

		if err := storage.UpdateAtom(atom); err != nil {
			return NoType, nil, err
		}

		return AtomsList, getAtomsList(), nil
	},

	AtomDelete: func(params *ActionParams) (ActionType, any, error) {
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return NoType, nil, ErrorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return NoType, nil, ErrorBadParams
		}

		if err := storage.DeleteAtom(id); err != nil {
			return NoType, nil, err
		}

		return AtomsList, getAtomsList(), nil
	},

	AtomCreate: func(params *ActionParams) (ActionType, any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return NoType, nil, ErrorBadParams
		}

		if atom.ID != nil || strings.TrimSpace(atom.Name) == "" {
			log.Println("error parsing params: bad atom")
			return NoType, nil, ErrorBadParams
		}

		storage.CreateAtom(atom)

		return AtomsList, getAtomsList(), nil
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
