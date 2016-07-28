package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type api struct {
	Codelabs   []codelab        `json:"codelabs"`
	Categories map[string]theme `json:"categories"`
}

func generateCodelabsAPI(codelabs []codelab, cats categories) (err error) {
	apiContent := api{
		Codelabs:   codelabs,
		Categories: cats.Categories,
	}

	var content []byte
	if content, err = json.MarshalIndent(apiContent, "", "  "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(apiPath, "codelab.json"), content, 0666); err != nil {
		return err
	}

	return nil
}
