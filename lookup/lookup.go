package lookup

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

// Entry represents a single entry in the CEDICT database
type Entry struct {
	Simplified  string
	Traditional string
	Pinyin      string
	Definitions []string
}

// Lookup table for CEDICT entries
type Lookup struct {
	Simplified  map[string][]*Entry
	Traditional map[string][]*Entry
}

func add(m map[string][]*Entry, k string, v *Entry) {
	if arr, ok := m[k]; ok {
		m[k] = append(arr, v)
	} else {
		m[k] = []*Entry{v}
	}
}

func canIgnore(token []byte) bool {
	return token == nil || len(token) == 0 || token[0] == '#'
}

// see https://cc-cedict.org/wiki/format:syntax
var pattern = regexp.MustCompile(`^(\S+) (\S+) \[(.+)\] /(.+)/$`)

// Build a lookup table, indexed by simplified and traditional
func Build(db io.Reader) (Lookup, error) {
	scanner := bufio.NewScanner(db)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if err == nil && !canIgnore(token) && !pattern.MatchString(string(token)) {
			err = errors.New("Cannot parse line: " + string(token))
		}
		return advance, token, err
	}
	scanner.Split(split)

	lookup := Lookup{make(map[string][]*Entry, 0), make(map[string][]*Entry, 0)}
	for scanner.Scan() {
		if canIgnore(scanner.Bytes()) {
			continue
		}
		match := pattern.FindAllStringSubmatch(scanner.Text(), 1)
		entry := Entry{
			Simplified:  match[0][2],
			Traditional: match[0][1],
			Pinyin:      match[0][3],
			Definitions: strings.Split(match[0][4], "/"),
		}
		add(lookup.Simplified, entry.Simplified, &entry)
		add(lookup.Traditional, entry.Traditional, &entry)
	}

	if err := scanner.Err(); err != nil {
		return Lookup{}, err
	}

	return lookup, nil
}
