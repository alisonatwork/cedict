package lookup

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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
	trie        *node
}

type node struct {
	word     bool
	children map[string]*node
}

// Match returns an array of the longest matching entries for the term
func (lookup *Lookup) Match(term string) []*Entry {
	s := strings.Split(term, "")
	matcher := ""
	match := ""
	matchidx := -1
	nomatch := ""
	n := lookup.trie
	ret := make([]*Entry, 0)
	for i := 0; i < len(s); i++ {
		matcher += s[i]
		n = n.children[s[i]]
		if n != nil && n.word {
			match = matcher
			matchidx = i
		}
		if n == nil || i == len(s)-1 {
			if matchidx >= 0 {
				if len(nomatch) > 0 {
					ret = append(ret, &Entry{Simplified: nomatch, Traditional: nomatch})
					nomatch = ""
				}
				entries := lookup.Simplified[match]
				if len(entries) == 0 {
					entries = lookup.Traditional[match]
				}
				ret = append(ret, entries...)
				i = matchidx
				match = ""
				matchidx = -1
			} else {
				nomatch += matcher
			}
			matcher = ""
			n = lookup.trie
		}
	}
	if len(nomatch) > 0 {
		ret = append(ret, &Entry{Simplified: nomatch, Traditional: nomatch})
	}
	return ret
}

func add(m map[string][]*Entry, k string, v *Entry) {
	if arr, ok := m[k]; ok {
		m[k] = append(arr, v)
	} else {
		m[k] = []*Entry{v}
	}
}

func addNode(lookup Lookup, entry *Entry) {
	s := strings.Split(entry.Simplified, "")
	t := strings.Split(entry.Traditional, "")
	if len(s) != len(t) {
		fmt.Fprintf(os.Stderr, "simplified and traditional not same length!? %s %s\n", entry.Simplified, entry.Traditional)
	}
	n := lookup.trie
	for i := range s {
		if n.children[s[i]] == nil {
			n.children[s[i]] = &node{
				word:     false,
				children: make(map[string]*node),
			}
			if s[i] != t[i] {
				n.children[t[i]] = n.children[s[i]]
			}
		}
		n = n.children[s[i]]
		if i == len(s)-1 {
			n.word = true
		}
	}
}

func canIgnore(token []byte) bool {
	return len(token) == 0 || token[0] == '#'
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

	lookup := Lookup{
		make(map[string][]*Entry),
		make(map[string][]*Entry),
		&node{
			word:     false,
			children: make(map[string]*node),
		},
	}
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
		addNode(lookup, &entry)
	}

	if err := scanner.Err(); err != nil {
		return Lookup{}, err
	}

	return lookup, nil
}
