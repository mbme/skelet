package storage

import "fmt"

type AtomType string

const (
	Record AtomType = ":record"
	File            = ":file"
)

type AtomId uint32

type AtomLink struct {
	Type *AtomType
	ID   *AtomId
}

func (l *AtomLink) IsValid() bool {
	return l.Type != nil && l.ID != nil
}

func (l *AtomLink) String() string {
	return fmt.Sprintf("%s/%v", l.Type, l.ID)
}

type Atom struct {
	Type AtomType
	ID   AtomId
	Name string
	Data string
}

func (a *Atom) String() string {
	return fmt.Sprintf("%s/%v[%s]", a.Type, a.ID, a.Name)
}
