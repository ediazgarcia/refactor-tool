package types

import (
	"encoding/json"
	"io/ioutil"
)

type Component struct {
	Name        string   `json:"name"`
	Props       []string `json:"props"`
	Hooks       []string `json:"hooks"`
	Imports     []string `json:"imports"`
	JSXElements []string `json:"jsxElements"`
	FilePath    string   `json:"filePath"`
	Content     string   `json:"content"`
}

type Report struct {
	Components        []ComponentReport `json:"components"`
	GlobalSuggestions []string          `json:"globalSuggestions"`
	Statistics        Statistics        `json:"statistics"`
}

type ComponentReport struct {
	Name        string           `json:"name"`
	File        string           `json:"file"`
	Suggestions []string         `json:"suggestions"`
	Metrics     ComponentMetrics `json:"metrics"`
}

type ComponentMetrics struct {
	Lines int `json:"lines"`
	Hooks int `json:"hooks"`
	Props int `json:"props"`
}

type Statistics struct {
	TotalComponents      int `json:"totalComponents"`
	ComponentsWithIssues int `json:"componentsWithIssues"`
	TotalSuggestions     int `json:"totalSuggestions"`
}

type RefactorConfig struct {
	MaxLines      int  `json:"maxLines"`
	MaxHooks      int  `json:"maxHooks"`
	MaxProps      int  `json:"maxProps"`
	UseMemo       bool `json:"useMemo"`
	UseArrowFuncs bool `json:"useArrowFuncs"`
	SortImports   bool `json:"sortImports"`
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(path string) (*RefactorConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &RefactorConfig{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
