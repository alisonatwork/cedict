package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hermanschaaf/cedict"
)

func openCeDict() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := home + "/tmp/cedict_ts.u8"

	return os.Open(path)
}

func main() {
	if len(os.Args) == 1 {
		return
	}

	f, err := openCeDict()
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Dictionary not found: getcedict first!\n")
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	lookup := make(map[string][]*cedict.Entry)
	for _, arg := range os.Args[1:] {
		lookup[arg] = make([]*cedict.Entry, 0, 1)
	}

	c := cedict.New(f)
	for err := c.NextEntry(); err == nil; err = c.NextEntry() {
		entry := c.Entry()
		if lookup[entry.Simplified] != nil {
			lookup[entry.Simplified] = append(lookup[entry.Simplified], entry)
		} else if lookup[entry.Traditional] != nil {
			lookup[entry.Traditional] = append(lookup[entry.Traditional], entry)
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		return
	}

	for _, v := range lookup {
		for _, e := range v {
			fmt.Printf("%s (%s) %s\n", e.Simplified, e.PinyinWithTones, strings.Join(e.Definitions, " / "))
		}
	}
}
