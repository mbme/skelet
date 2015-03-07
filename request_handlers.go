package main

import (
	"encoding/json"
	"errors"

	"log"

	s "github.com/mbme/skelet/storage"
)

var (
	errorNoHandler = errors.New("no handler for request")
	errorBadParams = errors.New("malformed request params")
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

var storage = s.NewStorage()

type atomInfo struct {
	ID         s.AtomID     `json:"id"`
	Type       s.AtomType   `json:"type"`
	Name       string       `json:"name"`
	Categories []s.Category `json:"categories"`
	TsCreated  *s.AtomTime  `json:"ts_created"`
	TsUpdated  *s.AtomTime  `json:"ts_updated"`
}

func toAtomInfo(atom *s.Atom) *atomInfo {
	return &atomInfo{
		ID:         *atom.ID,
		Type:       *atom.Type,
		Name:       atom.Name,
		Categories: atom.Categories,
		TsCreated:  atom.TsCreated,
		TsUpdated:  atom.TsUpdated,
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

var handlers = map[RequestMethod]func(*RequestParams) (any, error){
	AtomsListRead: func(_ *RequestParams) (any, error) {
		return getAtomsList(), nil
	},

	AtomRead: func(params *RequestParams) (any, error) {
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, errorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return nil, errorBadParams
		}

		atom, err := storage.GetAtom(id)
		if err != nil {
			log.Printf("can't find atom %s", id)
			return nil, err
		}

		return atom, nil
	},

	AtomCreate: func(params *RequestParams) (any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, errorBadParams
		}

		errors := atom.Validate()
		if atom.ID != nil {
			errors = append(errors, "id is not nil")
		}

		if atom.TsCreated != nil {
			errors = append(errors, "creation timestamp is not nil")
		}

		if atom.TsUpdated != nil {
			errors = append(errors, "creation timestamp is not nil")
		}

		if len(errors) > 0 {
			log.Printf("error parsing params: %v", errors)
			return errors, errorBadParams
		}

		storage.CreateAtom(atom)

		return atom, nil
	},

	AtomUpdate: func(params *RequestParams) (any, error) {
		atom := &s.Atom{}
		if err := params.readAs(atom); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, errorBadParams
		}

		errors := atom.Validate()
		if atom.ID == nil {
			errors = append(errors, "missing id")
		}

		if atom.TsCreated != nil {
			errors = append(errors, "creation timestamp is not nil")
		}

		if atom.TsUpdated != nil {
			errors = append(errors, "creation timestamp is not nil")
		}

		if len(errors) > 0 {
			log.Printf("error parsing params: %v", errors)
			return errors, errorBadParams
		}

		if err := storage.UpdateAtom(atom); err != nil {
			return nil, err
		}

		return atom, nil
	},

	AtomDelete: func(params *RequestParams) (any, error) {
		id := new(s.AtomID)
		if err := params.readAs(id); err != nil {
			log.Printf("error parsing params: %v", err)
			return nil, errorBadParams
		}

		if id == nil {
			log.Println("error parsing params: can't parse id")
			return nil, errorBadParams
		}

		if err := storage.DeleteAtom(id); err != nil {
			return nil, err
		}

		return id, nil
	},
}

// ProcessRequest handle client request
func ProcessRequest(reqType RequestMethod, params *RequestParams) (any, error) {
	handler, ok := handlers[reqType]

	if !ok {
		return nil, errorNoHandler
	}

	return handler(params)
}
