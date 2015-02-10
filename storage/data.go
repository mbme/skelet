package storage

import "errors"

type Storager interface {
	GetAtoms() []*Atom
	GetAtom(*AtomID) (*Atom, error)
	CreateAtom(*Atom)
	UpdateAtom(*Atom) error
	DeleteAtom(*AtomID) error
}

var (
	ErrorAtomNotFound = errors.New("atom not found")
)

func newRecord(id AtomID, name, data string) *Atom {
	atomType := Record

	return &Atom{
		Type: &atomType,
		ID:   &id,
		Name: name,
		Data: data,
	}
}

var records = map[AtomID]*Atom{}

type virtualStorage struct {
}

// NewStorage create new Storage instance
func NewStorage() Storager {
	for i, rec := range rawData {
		id := AtomID(i)
		records[id] = newRecord(id, rec.Name, rec.Data)
	}
	return &virtualStorage{}
}

func (l *virtualStorage) GetAtoms() []*Atom {
	var atoms []*Atom
	for _, a := range records {
		atoms = append(atoms, a)
	}

	return atoms
}

func (l *virtualStorage) GetAtom(id *AtomID) (*Atom, error) {
	atom, ok := records[*id]

	if !ok {
		return nil, ErrorAtomNotFound
	}

	return atom, nil
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

func (l *virtualStorage) DeleteAtom(id *AtomID) error {
	if _, err := l.GetAtom(id); err != nil {
		return err
	}

	delete(records, *id)

	return nil
}

func (l *virtualStorage) getNewId() *AtomID {
	maxID := AtomID(0)

	for id := range records {
		if id > maxID {
			maxID = id
		}
	}

	newID := maxID + 1

	return &newID
}

func (l *virtualStorage) CreateAtom(atom *Atom) {
	newID := l.getNewId()
	atom.ID = newID
	records[*newID] = atom
}
