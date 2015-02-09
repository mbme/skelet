package storage

import "fmt"

// AtomType is type of atom
type AtomType string

// possible atom types
const (
	Record AtomType = ":record"
	File            = ":file"
)

// AtomID is id of atom
type AtomID uint32

func (id *AtomID) String() string {
	return fmt.Sprintf("%v", *id)
}

// Atom is one information piece
type Atom struct {
	ID   *AtomID   `json:"id"`
	Type *AtomType `json:"type"`
	Name string    `json:"name"`
	Data string    `json:"data"`
}

func (a *Atom) String() string {
	return fmt.Sprintf("%v[%v/%s]", &a.ID, &a.Type, a.Name)
}

func (a *Atom) IsValid() bool {
	return a.ID != nil && a.Type != nil
}
