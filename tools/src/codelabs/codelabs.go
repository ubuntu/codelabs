package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

type codelab struct {
	Source     string   `json:"source"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Category   []string `json:"category"`
	Difficulty int      `json:"difficulty"`
	Duration   int      `json:"duration"`
	Tags       []string `json:"tags"`
	Updated    string   `json:"updated"`
	URL        string   `json:"url"`
}

func fetchAllCodelabs(codelabDir string) (codelabsMap []codelab, codelabsDir map[string]string, err error) {
	dirs, err := ioutil.ReadDir(codelabDir)
	if err != nil {
		return nil, nil, err
	}

	answers, errs, exit := make(chan codelab), make(chan error), make(chan bool)
	for _, dir := range dirs {
		// for each dir, extract metadata
		go func(dir os.FileInfo) {
			if !dir.IsDir() {
				exit <- true
				return
			}

			var thisCodelab codelab

			codelabjsonPath := path.Join(codelabPath, dir.Name(), metaFilename)
			content, erro := ioutil.ReadFile(codelabjsonPath)
			if erro != nil {
				errs <- erro
				return
			}

			if erro = json.Unmarshal(content, &thisCodelab); erro != nil {
				errs <- erro
				return
			}

			// send new codelab to channel
			answers <- thisCodelab

		}(dir)
	}

	codelabsMap = make([]codelab, 0)
	codelabsDir = make(map[string]string)
	var resultErrs bytes.Buffer

	for i := 0; i < len(dirs); i++ {
		select {
		case newCodelab := <-answers:
			codelabsMap = append(codelabsMap, newCodelab)
			codelabsDir[newCodelab.Source] = newCodelab.URL
		case err = <-errs:
			resultErrs.WriteString(err.Error())
		case _ = <-exit:
		}
	}
	close(answers)
	close(errs)
	close(exit)

	// error handling
	if resultErrs.String() != "" {
		return nil, nil, errors.New(resultErrs.String())
	}

	return codelabsMap, codelabsDir, err

}
