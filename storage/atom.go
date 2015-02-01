package storage

import "fmt"

type AtomType string
type AtomId uint32

const (
	Record AtomType = ":record"
	File            = ":file"
)

type Atom struct {
	Type AtomType
	ID   AtomId
	Name string
	Data string
}

func (a *Atom) String() string {
	return fmt.Sprintf("%s/%v[%s]", a.Type, a.ID, a.Name)
}
