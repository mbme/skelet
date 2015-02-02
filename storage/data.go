package storage

type Storager interface {
	GetAtoms() []*Atom
	GetAtom(*AtomLink) *Atom
}

var idGenerator AtomId

func newRecord(name, data string) *Atom {
	idGenerator++
	return &Atom{
		Type: Record,
		ID:   idGenerator,
		Name: name,
		Data: data,
	}
}

var records = []*Atom{
	newRecord("123", "fdas"),
	newRecord("124", "some text"),
}

type VirtualStorage struct {
}

func (l *VirtualStorage) GetAtoms() []*Atom {
	return records
}

func (l *VirtualStorage) GetAtom(link *AtomLink) *Atom {
	for _, atom := range l.GetAtoms() {
		if atom.Type == *link.Type && atom.ID == *link.ID {
			return atom
		}
	}
	return nil
}
