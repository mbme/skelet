package storage

import "errors"
import "time"

var errorAtomNotFound = errors.New("atom not found")

type virtualStorage struct {
}

func newRecord(id AtomID, name, data string, categories []Category) *Atom {
	atomType := Record

	now := AtomTime(time.Now())
	return &Atom{
		Type:       &atomType,
		ID:         &id,
		Name:       name,
		Data:       data,
		Categories: categories,
		TsCreated:  &now,
		TsUpdated:  &now,
	}
}

func (l *virtualStorage) getNewID() *AtomID {
	maxID := AtomID(0)

	for id := range records {
		if id > maxID {
			maxID = id
		}
	}

	newID := maxID + 1

	return &newID
}

var records = map[AtomID]*Atom{}

// NewStorage create new Storage instance
func NewStorage() Storager {
	for i, rec := range rawData {
		id := AtomID(i)
		records[id] = newRecord(id, rec.Name, rec.Data, rec.Categories)
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

func (l *virtualStorage) CreateAtom(atom *Atom) {
	now := AtomTime(time.Now())
	atom.TsCreated = &now
	atom.TsUpdated = &now

	newID := l.getNewID()
	atom.ID = newID

	records[*newID] = atom
}

func (l *virtualStorage) GetAtom(id *AtomID) (*Atom, error) {
	atom, ok := records[*id]

	if !ok {
		return nil, errorAtomNotFound
	}

	return atom, nil
}

func (l *virtualStorage) UpdateAtom(newAtom *Atom) error {
	atom, err := l.GetAtom(newAtom.ID)
	if err != nil {
		return err
	}

	now := AtomTime(time.Now())
	atom.TsUpdated = &now

	atom.Type = newAtom.Type
	atom.Name = newAtom.Name
	atom.Data = newAtom.Data
	atom.Categories = newAtom.Categories

	return nil
}

func (l *virtualStorage) DeleteAtom(id *AtomID) error {
	if _, err := l.GetAtom(id); err != nil {
		return err
	}

	delete(records, *id)

	return nil
}
