package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"refactor/types"
)

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *types.RefactorConfig {
	return &types.RefactorConfig{
		MaxLines:      300,
		MaxHooks:      5,
		MaxProps:      10,
		UseMemo:       true,
		UseArrowFuncs: true,
		SortImports:   true,
	}
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(path string) (*types.RefactorConfig, error) {
	// Si no existe el archivo, retornar config por defecto
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &types.RefactorConfig{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig guarda la configuración en un archivo JSON
func SaveConfig(config *types.RefactorConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
