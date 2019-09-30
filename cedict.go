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

func printEntries(defs []*lookup.Entry) {
	for _, e := range defs {
		fmt.Printf("%s\t[%s]\t/%s/\n", e.Simplified, pinyin.NumberToMark(e.Pinyin), pinyin.NumberToMark(strings.Join(e.Definitions, "/")))
	}
}

func output(strategy matchStrategy, lookup lookup.Lookup, word string) {
	defs := lookup.Simplified[word]
	if len(defs) == 0 {
		defs = lookup.Traditional[word]
	}
	if len(defs) == 0 {
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
			defs = lookup.Traditional[c]
			if len(defs) == 0 {
				fmt.Printf("%s", c)
			} else {
				fmt.Printf("%s", defs[0].Simplified)
			}
		}
		fmt.Printf("\n")
	} else {
		printEntries(defs)
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
