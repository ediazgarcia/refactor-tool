package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"refactor/analyzer"
	"refactor/refactor"
	"refactor/types"
)

var (
	projectPath  string
	configPath   string
	outputFormat string
	dryRun       bool
	verbose      bool
)

func init() {
	flag.StringVar(&projectPath, "path", "", "Path to React project")
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&outputFormat, "format", "json", "Output format (json|text)")
	flag.BoolVar(&dryRun, "dry-run", false, "Run analysis without applying changes")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s -path ./my-react-app -config ./config.json\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if projectPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Validar ruta del proyecto
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		log.Fatalf("Project path does not exist: %s", projectPath)
	}

	// Inicializar herramienta
	tool := NewRefactorTool(projectPath)

	// Cargar configuración si existe
	if configPath != "" {
		if err := tool.LoadConfig(configPath); err != nil {
			log.Printf("Warning: Could not load config file: %v", err)
		}
	}

	// Realizar análisis
	report := tool.AnalyzeProject()

	// Mostrar resultados
	if err := printReport(report); err != nil {
		log.Fatalf("Error printing report: %v", err)
	}

	// Aplicar refactorización si no es dry-run
	if !dryRun {
		if err := tool.ApplyRefactoring(); err != nil {
			log.Fatalf("Error applying refactoring: %v", err)
		}
		fmt.Println("\nRefactoring completed successfully!")
	}
}

type RefactorTool struct {
	projectPath string
	analyzer    *analyzer.Analyzer
	refactorer  *refactor.Refactorer
	components  []*types.Component
}

func NewRefactorTool(path string) *RefactorTool {
	return &RefactorTool{
		projectPath: path,
		analyzer:    analyzer.NewAnalyzer(),
		refactorer:  refactor.NewRefactorer(),
		components:  make([]*types.Component, 0),
	}
}

func (rt *RefactorTool) LoadConfig(path string) error {
	config, err := types.LoadConfig(path)
	if err != nil {
		return err
	}

	rt.analyzer.SetConfig(config)
	rt.refactorer.SetConfig(config)
	return nil
}

func (rt *RefactorTool) AnalyzeProject() *types.Report {
	if verbose {
		fmt.Println("Starting project analysis...")
	}

	report := &types.Report{
		Components:        make([]types.ComponentReport, 0),
		GlobalSuggestions: make([]string, 0),
		Statistics:        types.Statistics{},
	}

	// Encontrar archivos React
	files, err := rt.findReactFiles()
	if err != nil {
		log.Fatalf("Error finding React files: %v", err)
	}

	report.Statistics.TotalComponents = len(files)

	// Analizar cada componente
	for _, file := range files {
		if verbose {
			fmt.Printf("Analyzing component: %s\n", file)
		}

		component, err := rt.analyzer.ParseComponent(file)
		if err != nil {
			log.Printf("Error parsing component %s: %v", file, err)
			continue
		}

		rt.components = append(rt.components, component)

		componentReport := rt.analyzeComponent(component)
		report.Components = append(report.Components, componentReport)

		if len(componentReport.Suggestions) > 0 {
			report.Statistics.ComponentsWithIssues++
			report.Statistics.TotalSuggestions += len(componentReport.Suggestions)
		}
	}

	// Análisis global del proyecto
	report.GlobalSuggestions = rt.analyzeProjectStructure()

	return report
}

func (rt *RefactorTool) findReactFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(rt.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Ignorar node_modules y directorios ocultos
			if info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Comprobar extensiones de archivo
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jsx" || ext == ".tsx" || (ext == ".js" && !strings.Contains(path, ".test.")) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func (rt *RefactorTool) analyzeComponent(component *types.Component) types.ComponentReport {
	suggestions := rt.analyzer.AnalyzeComponent(component)

	return types.ComponentReport{
		Name:        component.Name,
		File:        component.FilePath,
		Suggestions: suggestions,
		Metrics: types.ComponentMetrics{
			Lines: len(strings.Split(component.Content, "\n")),
			Hooks: len(component.Hooks),
			Props: len(component.Props),
		},
	}
}

func (rt *RefactorTool) analyzeProjectStructure() []string {
	var suggestions []string

	// Analizar estructura de directorios
	suggestions = append(suggestions, rt.analyzer.AnalyzeProjectStructure(rt.projectPath)...)

	// Analizar dependencias comunes
	if deps := rt.analyzer.AnalyzeDependencies(rt.components); len(deps) > 0 {
		suggestions = append(suggestions, deps...)
	}

	return suggestions
}

func (rt *RefactorTool) ApplyRefactoring() error {
	if len(rt.components) == 0 {
		return fmt.Errorf("no components to refactor")
	}

	for _, component := range rt.components {
		if verbose {
			fmt.Printf("Refactoring component: %s\n", component.FilePath)
		}

		if err := rt.refactorer.ApplyRefactoring(component); err != nil {
			return fmt.Errorf("error refactoring %s: %w", component.FilePath, err)
		}
	}

	// Formatear código después de refactorizar
	if err := rt.refactorer.FormatCode(rt.projectPath); err != nil {
		return fmt.Errorf("error formatting code: %w", err)
	}

	return nil
}

func printReport(report *types.Report) error {
	switch outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
	case "text":
		fmt.Printf("\nAnalysis Report:\n")
		fmt.Printf("Total Components: %d\n", report.Statistics.TotalComponents)
		fmt.Printf("Components with Issues: %d\n", report.Statistics.ComponentsWithIssues)
		fmt.Printf("Total Suggestions: %d\n\n", report.Statistics.TotalSuggestions)

		if len(report.GlobalSuggestions) > 0 {
			fmt.Println("Global Suggestions:")
			for _, suggestion := range report.GlobalSuggestions {
				fmt.Printf("- %s\n", suggestion)
			}
			fmt.Println()
		}

		fmt.Println("Component Analysis:")
		for _, comp := range report.Components {
			if len(comp.Suggestions) > 0 {
				fmt.Printf("\n%s (%s):\n", comp.Name, comp.File)
				for _, suggestion := range comp.Suggestions {
					fmt.Printf("  - %s\n", suggestion)
				}
			}
		}
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}

	return nil
}
