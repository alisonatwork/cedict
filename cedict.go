package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alisonatwork/cedict/db"
	"github.com/alisonatwork/cedict/lookup"
	"github.com/alisonatwork/cedict/pinyin"
)

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

	var words []string
	if len(os.Args) > 2 && os.Args[1] == "-m" {
		words = append(strings.Split(os.Args[2], ""), os.Args[3:]...)
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
		defs := lookup.Simplified[word]
		if len(defs) == 0 {
			defs = lookup.Traditional[word]
		}
		if len(defs) == 0 {
			fmt.Printf("%s\n", word)
		} else {
			for _, e := range defs {
				fmt.Printf("%s\t[%s]\t/%s/\n", e.Simplified, pinyin.NumberToMark(e.Pinyin), pinyin.NumberToMark(strings.Join(e.Definitions, "/")))
			}
		}
	}
}
