package storage

// Storager general storage interface
type Storager interface {
	GetAtoms() []*Atom
	GetAtom(*AtomID) (*Atom, error)
	CreateAtom(*Atom)
	UpdateAtom(*Atom) error
	DeleteAtom(*AtomID) error
}
