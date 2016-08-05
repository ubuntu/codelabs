package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type api struct {
	Codelabs   []codelab        `json:"codelabs"`
	Categories map[string]theme `json:"categories"`
	Events     map[string]event `json:"events"`
}

func generateCodelabsAPI(codelabs []codelab, cats categoriesEvents) (err error) {
	apiContent := api{
		Codelabs:   codelabs,
		Categories: cats.Categories,
		Events:     cats.Events,
	}

	var content []byte
	if content, err = json.MarshalIndent(apiContent, "", "  "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(apiPath, "codelabs.json"), content, 0666); err != nil {
		return err
	}

	return nil
}
