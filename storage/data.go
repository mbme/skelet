package storage

type Storager interface {
	GetAtoms() []*Atom
	GetAtom(AtomID) *Atom
}

func newRecord(id int, name, data string) *Atom {
	return &Atom{
		Type: Record,
		ID:   AtomID(id),
		Name: name,
		Data: data,
	}
}

var records = []*Atom{}

type virtualStorage struct {
}

func NewStorage() Storager {
	for i, rec := range rawData {
		records = append(records, newRecord(i, rec.Name, rec.Data))
	}
	return &virtualStorage{}
}

func (l *virtualStorage) GetAtoms() []*Atom {
	return records
}

func (l *virtualStorage) GetAtom(id AtomID) *Atom {
	for _, atom := range l.GetAtoms() {
		if atom.ID == id {
			return atom
		}
	}
	return nil
}
