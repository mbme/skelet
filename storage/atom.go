package storage

import "fmt"

type AtomType string

const (
	Record AtomType = ":record"
	File            = ":file"
)

type Atom struct {
	Type AtomType
	ID   uint32
	Name string
	Data string
}

func (a *Atom) String() string {
	return fmt.Sprintf("%s/%v[%s]", a.Type, a.ID, a.Name)
}
