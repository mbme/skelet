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

// Atom is one information piece
type Atom struct {
	Type AtomType `json:"type"`
	ID   AtomID   `json:"id"`
	Name string   `json:"name"`
	Data string   `json:"data"`
}

func (a *Atom) String() string {
	return fmt.Sprintf("%s/%v[%s]", a.Type, a.ID, a.Name)
}
