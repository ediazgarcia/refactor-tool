package refactor

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"refactor/types"
)

type Refactorer struct {
	config *types.RefactorConfig
}

func NewRefactorer() *Refactorer {
	return &Refactorer{
		config: &types.RefactorConfig{
			MaxLines:      300,
			MaxHooks:      5,
			MaxProps:      10,
			UseMemo:       true,
			UseArrowFuncs: true,
			SortImports:   true,
		},
	}
}

// SetConfig establece la configuración del refactorizador
func (r *Refactorer) SetConfig(config *types.RefactorConfig) {
	r.config = config
}

func (r *Refactorer) ApplyRefactoring(component *types.Component) error {
	// Crear backup
	backupPath := component.FilePath + ".backup"
	if err := ioutil.WriteFile(backupPath, []byte(component.Content), 0644); err != nil {
		return fmt.Errorf("error creating backup: %w", err)
	}

	// Aplicar refactorización
	content := component.Content

	// Convertir a arrow function si es necesario
	if r.config.UseArrowFuncs && strings.Contains(content, "function ") {
		content = r.convertToArrowFunction(content, component.Name)
	}

	// Añadir memo si es necesario
	if r.config.UseMemo && len(component.Props) > 2 && !strings.Contains(content, "memo") {
		content = r.addMemo(content, component.Name)
	}

	// Ordenar imports
	if r.config.SortImports {
		content = r.sortImports(content, component.Imports)
	}

	// Escribir contenido refactorizado
	if err := ioutil.WriteFile(component.FilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing refactored content: %w", err)
	}

	return nil
}

func (r *Refactorer) convertToArrowFunction(content, componentName string) string {
	re := regexp.MustCompile(fmt.Sprintf(`function\s+%s\s*\((.*?)\)\s*{`, componentName))
	return re.ReplaceAllString(content, fmt.Sprintf(`const %s = ($1) => {`, componentName))
}

func (r *Refactorer) addMemo(content, componentName string) string {
	if !strings.Contains(content, "import { memo }") {
		content = "import { memo } from 'react';\n" + content
	}

	re := regexp.MustCompile(fmt.Sprintf(`export default %s`, componentName))
	return re.ReplaceAllString(content, fmt.Sprintf(`export default memo(%s)`, componentName))
}

func (r *Refactorer) sortImports(content string, imports []string) string {
	if len(imports) == 0 {
		return content
	}

	// Agrupar imports por tipo
	var reactImports, thirdPartyImports, localImports []string
	for _, imp := range imports {
		switch {
		case strings.Contains(imp, "react"):
			reactImports = append(reactImports, imp)
		case strings.HasPrefix(imp, "."):
			localImports = append(localImports, imp)
		default:
			thirdPartyImports = append(thirdPartyImports, imp)
		}
	}

	// Ordenar cada grupo
	sort.Strings(reactImports)
	sort.Strings(thirdPartyImports)
	sort.Strings(localImports)

	// Combinar todos los imports
	allImports := append(append(reactImports, thirdPartyImports...), localImports...)
	newImports := strings.Join(allImports, "\n")

	// Reemplazar sección de imports
	importSection := regexp.MustCompile(`(?s)(import.*?\n)+`)
	return importSection.ReplaceAllString(content, newImports+"\n\n")
}

func (r *Refactorer) FormatCode(projectPath string) error {
	cmd := exec.Command("npx", "prettier", "--write", ".")
	cmd.Dir = projectPath
	return cmd.Run()
}
