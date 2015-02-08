package storage

import "testing"

func TestVirtualStorage(t *testing.T) {
	s := NewStorage()

	if len(s.GetAtoms()) != len(records) {
		t.Errorf("virtual storage returns wrong records count: %v", len(records))
	}
}
