package main

import (
	"encoding/json"
	"errors"

	"log"

	"strings"

	s "github.com/mbme/skelet/storage"
)

//RequestMethod types of client requests
type RequestMethod string

//RequestParams raw parameters
type RequestParams json.RawMessage

//Possible requests
const (
	AtomsListRead RequestMethod = "atoms-list-read"

	AtomRead   = "atom-read"
	AtomCreate = "atom-create"
	AtomUpdate = "atom-update"
	AtomDelete = "atom-delete"

	NoType = ""
)

var (
	ErrorNoHandler = errors.New("no handler for request")
	ErrorBadParams = errors.New("malformed request params")
)

var storage = s.NewStorage()

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

type requestHandler func(*RequestParams) (any, error)

var handlers = map[RequestMethod]requestHandler{
	AtomsListRead: func(_ *RequestParams) (any, error) {
		return getAtomsList(), nil
	},

	AtomRead: func(params *RequestParams) (any, error) {
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, ErrorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return nil, ErrorBadParams
		}

		atom, err := storage.GetAtom(id)
		if err != nil {
			log.Printf("can't find atom %s", id)
			return nil, err
		}

		return atom, nil
	},

	AtomUpdate: func(params *RequestParams) (any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, ErrorBadParams
		}

		if atom.ID == nil || atom.Type == nil || !atom.Type.IsValid() {
			log.Println("error parsing params: bad atom")
			return nil, ErrorBadParams
		}

		if err := storage.UpdateAtom(atom); err != nil {
			return nil, err
		}

		return getAtomsList(), nil
	},

	AtomDelete: func(params *RequestParams) (any, error) {
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, ErrorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return nil, ErrorBadParams
		}

		if err := storage.DeleteAtom(id); err != nil {
			return nil, err
		}

		return getAtomsList(), nil
	},

	AtomCreate: func(params *RequestParams) (any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, ErrorBadParams
		}

		if atom.ID != nil || strings.TrimSpace(atom.Name) == "" {
			log.Println("error parsing params: bad atom")
			return nil, ErrorBadParams
		}

		storage.CreateAtom(atom)

		return getAtomsList(), nil
	},
}

// ProcessRequest handle client request
func ProcessRequest(reqType RequestMethod, params *RequestParams) (any, error) {
	handler, ok := handlers[reqType]

	if !ok {
		return nil, ErrorNoHandler
	}

	return handler(params)
}
