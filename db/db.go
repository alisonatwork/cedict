package db

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const fileName = "cedict_1_0_ts_utf-8_mdbg.txt"
const gzipUrl = "https://www.mdbg.net/chinese/export/cedict/" + fileName + ".gz"

func getLocalPath() (string, error) {
	appdata, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Join(appdata, "cedict"), os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Join(appdata, "cedict", fileName), nil
}

// Download the CC-CEDICT database from the web into a local cache
func Download() error {

	fmt.Printf("Connecting to %s ...\n", gzipUrl)

	resp, err := http.Get(gzipUrl)
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

// Open the local CC-CEDICT database
func Open() (*os.File, error) {
	path, err := getLocalPath()
	if err != nil {
		return nil, err
	}

	return os.Open(path)
}
