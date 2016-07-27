// Copyright 2016 Canonical
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	globalGA string
	version  = "0.1"
)

const (
	// metaFilename is codelab metadata file.
	metaFilename = "codelab.json"

	// claat tool executable
	claatFileName = "claat-linux-amd64"
	// FIXME: should be a proper dynamic link later on, with arch and other "latest" api
	claatURL  = "https://github.com/googlecodelabs/tools/releases/download/v0.4.0/" + claatFileName
	claatExec = "/tmp/" + claatFileName
)

var (
	// commands contains all valid subcommands, e.g. "codelabs-tool add".
	commands = map[string]func(){
		"add":     cmdAdd,
		"update":  cmdUpdate,
		"remove":  cmdRemove,
		"help":    usage,
		"version": func() { fmt.Println(version) },
	}
)

func cmdAdd() {
	if flag.NArg() == 0 {
		fatalf("Need at least one codelab to import. Try '-h' for options.")
	}

	if err := getClaat(); err != nil {
		fatalf("Couldn't download %s command: %v", claatURL, err)
	}

	ensureInCodelabDir()

	args := unique(flag.Args())
	printf(strings.Join(args, ", "))

	printf(globalGA)
}

func cmdUpdate() {
	if err := getClaat(); err != nil {
		fatalf("Couldn't download %s command: %v", claatURL, err)
	}

	ensureInCodelabDir()
}

func cmdRemove() {
	if flag.NArg() == 0 {
		fatalf("Need at least one codelab to remove. Try '-h' for options.")
	}

	ensureInCodelabDir()

}

// cd into codelab directory. Exit if failing
func ensureInCodelabDir() {
	codelabPath, err := getCodeLabDir()
	if err != nil {
		fatalf("Couldn't find codelab directory: %v", err)
	}
	os.Chdir(codelabPath)
}

// download or reuse existing claat binary from temp dir
func getClaat() (err error) {
	if _, err := os.Stat(claatExec); err != nil {
		printf("Downloading claat tool")
		// Create the file
		out, err := os.Create(claatExec)
		if err != nil {
			return err
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(claatURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	os.Chmod(claatExec, 755)

	return nil
}

// find codelab dir path relative to current executable
func getCodeLabDir() (codelabDir string, err error) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return codelabDir, err
	}

	for codelabDir == "" {
		dir = path.Clean(path.Join(dir, ".."))
		if dir == "/" {
			return "", errors.New("Couldn't find any codelab directory")
		}
		codelabDir = path.Join(dir, "src", "codelabs")
		_, err := os.Stat(codelabDir)
		_, err2 := os.Stat(path.Join(dir, "tools"))
		if err != nil || err2 != nil {
			codelabDir = ""
		}
	}

	return codelabDir, nil
}

// printf prints formatted string fmt with args to stderr.
func printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// fatalf calls printf and exits immediatly with non-zero code.
func fatalf(format string, args ...interface{}) {
	printf(format, args...)
	os.Exit(1)
}

// unique de-dupes a.
// The argument a is not modified.
func unique(a []string) []string {
	seen := make(map[string]struct{}, len(a))
	res := make([]string, 0, len(a))
	for _, s := range a {
		if _, y := seen[s]; !y {
			res = append(res, s)
			seen[s] = struct{}{}
		}
	}
	return res
}

func main() {
	log.SetFlags(0)
	flag.StringVar(&globalGA, "ga", "UA-81281030-1", "global Google Analytics account")

	if len(os.Args) == 1 {
		fatalf("Need subcommand. Try '-h' for options.")
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		return
	}

	cmd := commands[os.Args[1]]
	if cmd == nil {
		fatalf("Unknown subcommand. Try '-h' for options.")
	}
	flag.Usage = usage
	flag.CommandLine.Parse(os.Args[2:])
	cmd()
	os.Exit(0)
}

// usage prints usageText and program arguments to stderr.
func usage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `Usage: codelab-tool <cmd>

where cmd can :
* add [flags] google_id [google_id…]
* update [flags]
* rm codelab_dirname|google_id [codelab_dirname|google_id…]
* version

Add and update flags are:
- ga to specify the global GA account

## Add command

Add takes one or more Google doc ID as arguments (omitting https://docs.google.com/... part)
and import it as a new codelabs.

The program exits with non-zero code if at least one doc could not be added.

A new commit and tag is created on success.

## Update command

Update scans one or more local codelab for codelab.json metadata
files, recursively. A directory containing the metadata file is expected
to be a codelab previously created with the add command.

Directory is detected relative to this executable.

Unused codelab assets will be deleted, as well as the entire codelab directory,
if codelab ID has changed since last update or add (no more matching document).

In the latter case, where codelab ID has changed, the new directory
will be placed alongside the old one. In other words, it will have the same ancestor
as the old one.

The program does not follow symbolic links and exits with non-zero code
if no metadata found or at least one codelab could not be updated.

A new commit and tag is created on success.

## Remove command

Remove looks for codelabs matching 'src' dir name or codelab id to remove them from the
codelab list.

A new commit and tag is created on success.

`
