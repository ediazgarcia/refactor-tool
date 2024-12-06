.PHONY: build run test clean

# Variables
BINARY_NAME=refactor
GO_FILES=$(shell find . -name '*.go')

# Compilar el proyecto
build:
	go build -o bin/$(BINARY_NAME)

# Ejecutar el proyecto
run:
	go run main.go

# Ejecutar tests
test:
	go test ./...

# Limpiar archivos generados
clean:
	rm -rf bin/
	go clean

# Instalar dependencias
deps:
	go mod download
	go mod tidy

# Formatear código
fmt:
	go fmt ./...

# Verificar código
lint:
	go vet ./...

# Instalación completa
install: deps build

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make build    - Compila el proyecto"
	@echo "  make run      - Ejecuta el proyecto"
	@echo "  make test     - Ejecuta los tests"
	@echo "  make clean    - Limpia archivos generados"
	@echo "  make deps     - Instala dependencias"
	@echo "  make fmt      - Formatea el código"
	@echo "  make lint     - Verifica el código"
	@echo "  make install  - Instalación completa"