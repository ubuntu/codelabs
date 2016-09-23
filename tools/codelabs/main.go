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
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	globalGA    string
	codelabPath string
	apiPath     string
	version     = "0.1"
	catEvents   *categoriesEvents
)

const (
	// metaFilename is codelab metadata file.
	metaFilename = "codelab.json"

	// claat tool executable
	claatFileName = "claat-linux-amd64"
	// FIXME: we are using our fork for now due to "difficulty" additional tag
	claatURL  = "https://people.canonical.com/~didrocks/" + claatFileName
	claatExec = "./" + claatFileName
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

	args := unique(flag.Args())
	cmd := exec.Command(claatExec, "export", "-ga", globalGA, "-f", "ubuntu-template.html", "-o", codelabPath, "--prefix", "../../..", strings.Join(args, ", "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatalf("Couldn't add new codelab")
	}
}

func cmdUpdate() {
	if err := getClaat(); err != nil {
		fatalf("Couldn't download %s command: %v", claatURL, err)
	}

	cmd := exec.Command(claatExec, "update", "-ga", globalGA, "--prefix", "../../..", codelabPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatalf("Couldn't add refresh codelabs")
	}
}

func cmdRemove() {
	if flag.NArg() == 0 {
		fatalf("Need at least one codelab to remove. Try '-h' for options.")
	}

	_, codelabsIDToDir, err := fetchAllCodelabs(codelabPath)
	if err != nil {
		fatalf("Couldn't introspect existing codelabs: %s", err)
	}

	var failure bool
	for _, dir := range flag.Args() {
		if err := removeDir(path.Join(codelabPath, dir)); err != nil {
			// try to see if an ID was provided
			if _, present := codelabsIDToDir[dir]; present {
				if err2 := removeDir(path.Join(codelabPath, codelabsIDToDir[dir])); err2 != nil {
					log.Printf("Couldn't find or remove: %s", path.Join(codelabPath, codelabsIDToDir[dir]))
					failure = true
				}
			} else {
				log.Printf("Couldn't find or remove: %s", path.Join(codelabPath, dir))
				failure = true
			}
		}
	}

	if failure {
		printf("One or more codelabs couldn't get removed")
		// we don't exit here to regenerate current codelabs list
	}

}

func removeDir(path string) (err error) {
	if _, err = os.Stat(path); err != nil {
		return err
	}
	if err = os.RemoveAll(path); err != nil {
		log.Printf("Found, but couldn't remove %s: %v", path, err)
	}
	return err
}

// cd into codelab directory. Exit if failing
func ensureInToolsDir() {
	var toolsPath string
	var err error
	toolsPath, codelabPath, apiPath, err = getDirs()
	if err != nil {
		fatalf("Couldn't find tools or codelab directory: %v", err)
	}
	os.Chdir(toolsPath)
}

// download or reuse existing claat binary from temp dir
func getClaat() (err error) {
	if _, err := os.Stat(claatFileName); err != nil {
		printf("Downloading claat tool")
		// Create the file
		out, err := os.Create(claatFileName)
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

	os.Chmod(claatFileName, 0755)

	return nil
}

// find codelab and tools dir path relative to current executable
func getDirs() (toolsDir string, codelabDir string, apiPath string, err error) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return toolsDir, codelabDir, apiPath, err
	}

	rootDirFound := false
	for rootDirFound != true {
		if dir == "/" {
			return "", "", "", errors.New("Couldn't find any codelab or tools directory")
		}
		toolsDir = path.Join(dir, "tools")
		codelabDir = path.Join(dir, "src", "codelabs")
		apiPath = path.Join(dir, "api")
		_, toolsExistErr := os.Stat(toolsDir)
		_, codeLabsExistErr := os.Stat(codelabDir)
		_, rootExistErr := os.Stat(path.Join(dir, "bower.json"))
		if toolsExistErr == nil && codeLabsExistErr == nil && rootExistErr == nil {
			rootDirFound = true
		}
		dir = path.Clean(path.Join(dir, ".."))
	}

	return toolsDir, codelabDir, apiPath, nil
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
	flag.StringVar(&globalGA, "ga", "UA-1018242-64", "global Google Analytics account")

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

	ensureInToolsDir()
	var err error
	if catEvents, err = loadCategoriesData("../categories-events.json"); err != nil {
		fatalf("Couldn't load categories.json file: %s", err)
	}

	cmd()

	var codelabsMap []codelab
	codelabsMap, _, err = fetchAllCodelabs(codelabPath)
	if err != nil {
		fatalf("Couldn't introspect existing codelabs: %s", err)
	}

	if err = generateCodelabsAPI(codelabsMap, *catEvents); err != nil {
		fatalf("Couldn't save new categories.json api file: %s", err)
	}

	os.Exit(0)
}

// usage prints usageText and program arguments to stderr.
func usage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `Usage: codelab <cmd>

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
