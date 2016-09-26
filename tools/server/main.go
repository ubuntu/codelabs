package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func getRootDir() (rootDir string, err error) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return rootDir, err
	}

	for {
		if dir == "/" {
			return "", errors.New("Couldn't find root directory")
		}

		if _, rootExistErr := os.Stat(path.Join(dir, "src")); rootExistErr == nil {
			return dir, nil
		}

		dir = path.Clean(path.Join(dir, ".."))
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
	port := flag.Int("p", 8123, "Port to listen at")
	flag.Parse()

	rootDir, err := getRootDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(rootDir)

	http.HandleFunc("/", rootHandler)

	log.Printf("Listening on http://localhost:%d", *port)

	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
}
