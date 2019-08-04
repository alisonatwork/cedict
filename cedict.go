package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hermanschaaf/cedict"
)

const cedictGzipUrl = "https://www.mdbg.net/chinese/export/cedict/cedict_1_0_ts_utf-8_mdbg.txt.gz"

func getLocalPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/tmp/cedict_1_0_ts_utf-8_mdbg.txt", nil
}

func downloadCeDict() error {

	resp, err := http.Get(cedictGzipUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer r.Close()

	path, err := getLocalPath()
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Printf("Downloading: %s to %s\n", cedictGzipUrl, path)

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}
	return err
}

func openCeDict() (*os.File, error) {
	path, err := getLocalPath()
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}

func main() {
	if len(os.Args) == 1 {
		return
	}

	if len(os.Args) == 2 && os.Args[1] == "get" {
		err := downloadCeDict()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s!\n", err)
		}
		return
	}

	f, err := openCeDict()
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Dictionary not found: get first!\n")
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
