package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alisonatwork/cedict/db"
	"github.com/hermanschaaf/cedict"
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

	lookup := make(map[string][]*cedict.Entry)
	for _, word := range words {
		lookup[word] = make([]*cedict.Entry, 0, 1)
	}

	c := cedict.New(db)
	for err := c.NextEntry(); err == nil; err = c.NextEntry() {
		entry := c.Entry()
		if lookup[entry.Simplified] != nil {
			lookup[entry.Simplified] = append(lookup[entry.Simplified], entry)
		} else if lookup[entry.Traditional] != nil {
			lookup[entry.Traditional] = append(lookup[entry.Traditional], entry)
		}
	}

	err = db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	for _, word := range words {
		if len(lookup[word]) == 0 {
			fmt.Printf("%s\n", word)
		} else {
			for _, e := range lookup[word] {
				fmt.Printf("%s\t[%s]\t/%s/\n", e.Simplified, e.PinyinWithTones, strings.Join(e.Definitions, "/"))
			}
		}
	}
}
