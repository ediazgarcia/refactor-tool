package analyzer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"refactor/types"
)

type Analyzer struct {
	config         *types.RefactorConfig
	componentRegex *regexp.Regexp
	importRegex    *regexp.Regexp
	hookRegex      *regexp.Regexp
	propsRegex     *regexp.Regexp
	jsxRegex       *regexp.Regexp
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		config: &types.RefactorConfig{
			MaxLines:      300,
			MaxHooks:      5,
			MaxProps:      10,
			UseMemo:       true,
			UseArrowFuncs: true,
			SortImports:   true,
		},
		componentRegex: regexp.MustCompile(`function\s+(\w+)|class\s+(\w+)`),
		importRegex:    regexp.MustCompile(`import.*from.*`),
		hookRegex:      regexp.MustCompile(`use\w+\(`),
		propsRegex:     regexp.MustCompile(`\{(\w+)\}|\s(\w+)=`),
		jsxRegex:       regexp.MustCompile(`<(\w+)[^>]*>`),
	}
}

// SetConfig establece la configuración del analizador
func (a *Analyzer) SetConfig(config *types.RefactorConfig) {
	a.config = config
}

func (a *Analyzer) ParseComponent(filePath string) (*types.Component, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	strContent := string(content)

	// Extract component name
	matches := a.componentRegex.FindStringSubmatch(strContent)
	name := "Unknown"
	if len(matches) > 1 {
		if matches[1] != "" {
			name = matches[1]
		} else if matches[2] != "" {
			name = matches[2]
		}
	}

	return &types.Component{
		Name:        name,
		Props:       a.extractProps(strContent),
		Hooks:       a.extractHooks(strContent),
		Imports:     a.extractImports(strContent),
		JSXElements: a.extractJSX(strContent),
		FilePath:    filePath,
		Content:     strContent,
	}, nil
}

func (a *Analyzer) AnalyzeComponent(component *types.Component) []string {
	var suggestions []string

	// Verificar tamaño del componente
	lines := len(strings.Split(component.Content, "\n"))
	if lines > a.config.MaxLines {
		suggestions = append(suggestions,
			fmt.Sprintf("Component is too large (%d lines). Consider splitting it", lines))
	}

	// Verificar número de hooks
	if len(component.Hooks) > a.config.MaxHooks {
		suggestions = append(suggestions,
			fmt.Sprintf("Too many hooks (%d). Consider custom hooks", len(component.Hooks)))
	}

	// Verificar número de props
	if len(component.Props) > a.config.MaxProps {
		suggestions = append(suggestions,
			fmt.Sprintf("Too many props (%d). Consider composition", len(component.Props)))
	}

	return suggestions
}

// AnalyzeProjectStructure analiza la estructura del proyecto
func (a *Analyzer) AnalyzeProjectStructure(projectPath string) []string {
	var suggestions []string

	// Verificar estructura de directorios común en React
	expectedDirs := []string{"components", "pages", "hooks", "utils", "services"}
	for _, dir := range expectedDirs {
		if _, err := os.Stat(filepath.Join(projectPath, "src", dir)); os.IsNotExist(err) {
			suggestions = append(suggestions,
				fmt.Sprintf("Consider adding a '%s' directory for better organization", dir))
		}
	}

	return suggestions
}

// AnalyzeDependencies analiza las dependencias comunes entre componentes
func (a *Analyzer) AnalyzeDependencies(components []*types.Component) []string {
	var suggestions []string
	imports := make(map[string]int)

	// Contar imports comunes
	for _, comp := range components {
		for _, imp := range comp.Imports {
			imports[imp]++
		}
	}

	// Sugerir hooks personalizados para imports muy comunes
	for imp, count := range imports {
		if count > len(components)/2 && strings.Contains(imp, "useState") {
			suggestions = append(suggestions,
				fmt.Sprintf("Consider creating a custom hook for commonly used state logic (%s)", imp))
		}
	}

	return suggestions
}

func (a *Analyzer) extractProps(content string) []string {
	matches := a.propsRegex.FindAllStringSubmatch(content, -1)
	props := make([]string, 0)
	for _, match := range matches {
		if match[1] != "" {
			props = append(props, match[1])
		} else if match[2] != "" {
			props = append(props, match[2])
		}
	}
	return removeDuplicates(props)
}

func (a *Analyzer) extractHooks(content string) []string {
	return a.hookRegex.FindAllString(content, -1)
}

func (a *Analyzer) extractImports(content string) []string {
	return a.importRegex.FindAllString(content, -1)
}

func (a *Analyzer) extractJSX(content string) []string {
	return a.jsxRegex.FindAllString(content, -1)
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, entry := range slice {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
