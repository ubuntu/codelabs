package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getRootDir() (rootDir string, err error) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return rootDir, err
	}

	for {
		dir = path.Clean(path.Join(dir, ".."))
		if dir == "/" {
			return "", errors.New("Couldn't find root directory")
		}

		if _, rootExistErr := os.Stat(path.Join(dir, "bower.json")); rootExistErr == nil {
			return dir, nil
		}
	}

}

// serve file directly, unless it's an asset (containing "." in last world, if so, redirect to index.html)
func rootHandler(w http.ResponseWriter, req *http.Request) {
	splittedURL := strings.Split(req.URL.Path, "/")
	if strings.Contains(splittedURL[len(splittedURL)-1], ".") {
		http.ServeFile(w, req, req.URL.Path[1:])
	} else {
		http.ServeFile(w, req, "index.html")
	}
}

func main() {
	rootDir, err := getRootDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(rootDir)

	http.HandleFunc("/", rootHandler)

	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
}
