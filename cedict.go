package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hermanschaaf/cedict"
)

const cedictFile = "cedict_1_0_ts_utf-8_mdbg.txt"
const cedictGzipUrl = "https://www.mdbg.net/chinese/export/cedict/" + cedictFile + ".gz"

func getLocalPath() (string, error) {
	appdata, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(appdata, "cedict"), os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Join(appdata, "cedict", cedictFile), nil
}

func downloadCeDict() error {

	fmt.Printf("Connecting to %s ...\n", cedictGzipUrl)

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

	fmt.Printf("Downloading to %s ...\n", path)

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	fmt.Printf("Done!\n")
	return nil
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

	var words []string
	if len(os.Args) > 2 && os.Args[1] == "-m" {
		words = append(strings.Split(os.Args[2], ""), os.Args[3:]...)
	} else {
		words = os.Args[1:]
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
	for _, word := range words {
		lookup[word] = make([]*cedict.Entry, 0, 1)
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

	for _, word := range words {
		if len(lookup[word]) == 0 {
			fmt.Printf("%s\n", word)
		} else {
			for _, e := range lookup[word] {
				fmt.Printf("%s (%s) %s\n", e.Simplified, e.PinyinWithTones, strings.Join(e.Definitions, " / "))
			}
		}
	}
}
