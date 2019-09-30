package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alisonatwork/cedict/db"
	"github.com/alisonatwork/cedict/lookup"
	"github.com/alisonatwork/cedict/pinyin"
)

type matchStrategy int

const (
	exact matchStrategy = iota
	splitChar
)

func printEntries(entries []*lookup.Entry) {
	for _, e := range entries {
		defs := strings.Join(e.Definitions, "/")
		fmt.Printf("%s\t[%s]\t/%s/\n", e.Simplified, pinyin.NumberToMark(e.Pinyin), pinyin.NumberToMark(defs))
	}
}

func output(strategy matchStrategy, lookup lookup.Lookup, word string) {
	// we use simplified as the authoritative index lookup because it returns more hits
	// the problem: if we search 面 we want both noodles and face, not just face (simplified wins)
	//              if we search 碗 we just want bowl, not 4 different variants of bowl (traditional wins)
	entries := lookup.Simplified[word]
	if len(entries) == 0 {
		entries = lookup.Traditional[word]
	}
	if len(entries) == 0 {
		if len(word) > 1 {
			switch strategy {
			case exact:
				break
			case splitChar:
				for _, c := range strings.Split(word, "") {
					output(strategy, lookup, c)
				}
				return
			}
		}
		// convert traditional to simplified before printing unknown char
		for _, c := range strings.Split(word, "") {
			entries = lookup.Traditional[c]
			if len(entries) == 0 {
				fmt.Printf("%s", c)
			} else {
				fmt.Printf("%s", entries[0].Simplified)
			}
		}
		fmt.Printf("\n")
	} else {
		printEntries(entries)
	}
}

func main() {
	if len(os.Args) == 1 {
		return
	}

	if len(os.Args) == 2 && os.Args[1] == "get" {
		err := db.Download()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		}
		return
	}

	var strategy = exact
	var words []string
	if len(os.Args) > 2 && os.Args[1] == "-m" {
		words = os.Args[2:]
		strategy = splitChar
	} else {
		words = os.Args[1:]
	}

	db, err := db.Open()
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Dictionary not found: get first!\n")
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	lookup, err := lookup.Build(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	err = db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	for _, word := range words {
		output(strategy, lookup, word)
	}
}
