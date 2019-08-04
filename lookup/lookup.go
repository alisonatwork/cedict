package lookup

import (
	"os"

	"github.com/hermanschaaf/cedict"
)

// Lookup table for CEDICT entries
type Lookup struct {
	Simplified  map[string][]*cedict.Entry
	Traditional map[string][]*cedict.Entry
}

func add(m map[string][]*cedict.Entry, k string, v *cedict.Entry) {
	if arr, ok := m[k]; ok {
		m[k] = append(arr, v)
	} else {
		m[k] = []*cedict.Entry{v}
	}
}

// Build a lookup table, indexed by simplified and traditional
func Build(db *os.File) Lookup {
	lookup := Lookup{make(map[string][]*cedict.Entry, 0), make(map[string][]*cedict.Entry, 0)}
	cedict := cedict.New(db)
	for err := cedict.NextEntry(); err == nil; err = cedict.NextEntry() {
		entry := cedict.Entry()
		add(lookup.Simplified, entry.Simplified, entry)
		add(lookup.Traditional, entry.Traditional, entry)
	}
	return lookup
}
