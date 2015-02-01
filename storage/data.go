package storage

type Storager interface {
	GetAtoms() []*Atom
}

var idGenerator uint32

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
