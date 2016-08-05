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

type event struct {
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
}

type categoriesEvents struct {
	Categories map[string]theme `json:"categories"`
	Events     map[string]event `json:"events"`
}

func loadCategoriesData(filePath string) (c *categoriesEvents, err error) {
	categoriesFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(categoriesFileContent, &c); err != nil {
		return nil, err
	}

	return c, nil
}
