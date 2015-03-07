package storage

import "fmt"
import "strings"
import "time"
import "strconv"

// AtomType is type of atom
type AtomType string

// possible atom types
const (
	Record AtomType = ":record"
	File            = ":file"
)

func (t *AtomType) isValid() bool {
	return *t == Record || *t == File
}

// Category is atom category (context)
type Category string

// valid categories are trimmed lower-case strings
func (c Category) isValid() bool {
	str := string(c)
	trimmed := strings.TrimSpace(str)
	lower := strings.ToLower(str)
	return str == trimmed && str == lower
}

// AtomID is id of atom
type AtomID uint32

func (id *AtomID) String() string {
	return fmt.Sprintf("%v", *id)
}

// AtomTime timestamp
type AtomTime time.Time

// MarshalJSON implements json_Marshaller
func (t AtomTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// Atom is one information piece
type Atom struct {
	ID         *AtomID    `json:"id"`
	Type       *AtomType  `json:"type"`
	Name       string     `json:"name"`
	Data       string     `json:"data"`
	Categories []Category `json:"categories"`
	TsCreated  *AtomTime  `json:"ts_created"`
	TsUpdated  *AtomTime  `json:"ts_updated"`
}

func (a *Atom) String() string {
	return fmt.Sprintf("%v%v/%s %v", &a.ID, &a.Type, a.Name, a.Categories)
}

func unique(arr []Category) []Category {
	var result []Category
	seen := map[Category]int{}

	for _, category := range arr {
		if _, ok := seen[category]; !ok {
			result = append(result, category)
			seen[category] = 1
		}
	}

	return result
}

// Validate validates atom and returns array of errors
func (a *Atom) Validate() []string {
	var errors []string

	if a.Type == nil {
		errors = append(errors, "missing type")
	} else if !a.Type.isValid() {
		errors = append(errors, fmt.Sprintf("bad type: %v", a.Type))
	}

	if strings.TrimSpace(a.Name) == "" {
		errors = append(errors, "empty name")
	}

	if len(a.Categories) == 0 {
		errors = append(errors, "no categories specified")
	}

	hasBadCategory := false
	for _, category := range a.Categories {
		if !category.isValid() {
			hasBadCategory = true
			errors = append(errors, fmt.Sprintf("malformed category %s", category))
		}
	}

	if !hasBadCategory && len(a.Categories) != len(unique(a.Categories)) {
		errors = append(errors, "found duplicate categories")
	}

	return errors
}
