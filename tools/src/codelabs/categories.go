package main

import (
	"encoding/json"
	"io/ioutil"
)

type theme struct {
	MainColor      string `json:"maincolor"`
	SecondaryColor string `json:"secondarycolor"`
	LightColor     string `json:"lightcolor"`
}

type categories struct {
	Categories map[string]theme `json:"categories"`
}

func loadCategoriesData(filePath string) (categoriesData *categories, err error) {
	categoriesFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(categoriesFileContent, &categoriesData); err != nil {
		return nil, err
	}

	return categoriesData, nil
}
