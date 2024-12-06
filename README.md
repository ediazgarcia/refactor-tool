# React Code Refactor Tool

## 📋 Descripción

Una herramienta de línea de comandos desarrollada en Go para analizar y refactorizar automáticamente proyectos React/TypeScript. Identifica problemas comunes y propone mejoras siguiendo las mejores prácticas.

## ✨ Características

- Análisis automático de componentes React
- Detección de problemas de rendimiento
- Refactorización automática de código
- Generación de reportes detallados
- Soporte para TypeScript
- Formateo automático de código

## 🔍 Análisis

- Tamaño de componentes
- Número de hooks
- Cantidad de props
- Estructura del proyecto
- Patrones de código comunes
- Dependencias entre componentes

## 🛠️ Mejores Prácticas Implementadas

- Conversión a arrow functions
- Implementación de React.memo
- Optimización de hooks
- Ordenamiento de imports
- División de componentes grandes
- Extracción de custom hooks

## 🚀 Instalación

### Prerrequisitos

- Go 1.21 o superior
- Node.js y npm (para Prettier)

### Pasos de Instalación

```bash
# Clonar el repositorio
git clone https://github.com/ediazgarcia/refactor-tool.git
cd refactor-tool

# Instalar dependencias
go mod download
go mod tidy

# Compilar
go build -o refactorapp
```

## 📖 Uso

### Comando Básico

```bash
./refactorapp -path /ruta/a/tu/proyecto-react
```

### Opciones Disponibles

```bash
Options:
  -path string    Path to React project
  -config string  Path to config file (opcional)
  -format string  Output format (json|text) (default "json")
  -dry-run        Run analysis without applying changes
  -verbose        Show detailed output
```

## 📊 Ejemplo de Reporte

```json
{
  "components": [
    {
      "name": "ExampleComponent",
      "file": "src/components/Example.tsx",
      "suggestions": [
        "Component is too large (350 lines). Consider splitting it",
        "Too many hooks (7). Consider custom hooks"
      ],
      "metrics": {
        "lines": 350,
        "hooks": 7,
        "props": 12
      }
    }
  ],
  "globalSuggestions": [
    "Consider adding a 'utils' directory for better organization"
  ],
  "statistics": {
    "totalComponents": 25,
    "componentsWithIssues": 5,
    "totalSuggestions": 8
  }
}
```

## 🔧 Estructura del Proyecto

```refactor-tool/
├── main.go
├── analyzer/
│   └── analyzer.go
├── refactor/
│   └── refactor.go
├── types/
│   └── types.go
└── config/
    └── config.go
```

## ⚙️ Personalización

La herramienta puede personalizarse mediante un archivo de configuración que permite ajustar:

- Número máximo de hooks
- Número máximo de props
- Reglas de refactorización
- Formateo de código
- Límites líneas código

## 🤝 Contribución

Las contribuciones son bienvenidas. Por favor:

1. Fork el repositorio
2. Crea una rama para tu feature
3. Realiza tus cambios
4. Envía un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia

## ✍️ Autor

Desarrollado por Eddy Díaz.
