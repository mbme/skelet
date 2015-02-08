package storage

import "errors"

type Storager interface {
	GetAtoms() []*Atom
	GetAtom(*AtomID) (*Atom, error)
	UpdateAtom(*Atom) error
}

var (
	ErrorAtomNotFound = errors.New("atom not found")
)

func newRecord(id int, name, data string) *Atom {
	atomID := AtomID(id)
	atomType := Record

	return &Atom{
		Type: &atomType,
		ID:   &atomID,
		Name: name,
		Data: data,
	}
}

var records = []*Atom{}

type virtualStorage struct {
}

// NewStorage create new Storage instance
func NewStorage() Storager {
	for i, rec := range rawData {
		records = append(records, newRecord(i, rec.Name, rec.Data))
	}
	return &virtualStorage{}
}

func (l *virtualStorage) GetAtoms() []*Atom {
	return records
}

func (l *virtualStorage) GetAtom(id *AtomID) (*Atom, error) {
	for _, atom := range l.GetAtoms() {
		if atom.ID == id {
			return atom, nil
		}
	}
	return nil, ErrorAtomNotFound
}

func (l *virtualStorage) UpdateAtom(newAtom *Atom) error {
	atom, err := l.GetAtom(newAtom.ID)
	if err != nil {
		return err
	}

	atom.Type = newAtom.Type
	atom.Name = newAtom.Name
	atom.Data = newAtom.Data

	return nil
}
