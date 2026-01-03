package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// display displays the generated files/content
func display(outDir string) {
	printBar()
	log.Println("Generated Model Files:")

	pattern := filepath.Join(getModelOutDir(outDir), "*.go")
	modelFiles, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Error listing model files: %v", err)
	}

	for _, file := range modelFiles {
		displayFile(file)
	}
}

// displayFile reads and displays the content of a generated file
func displayFile(path string) {
	r, err := os.OpenRoot(os.TempDir())
	if err != nil {
		log.Fatalf("Error opening root directory: %v", err)
	}
	defer r.Close()

	path = strings.TrimPrefix(path, os.TempDir()+string(os.PathSeparator))
	content, err := r.ReadFile(path)
	if err != nil {
		log.Printf("Error reading %s: %v", path, err)
		return
	}

	log.Printf("\n--- %s ---\n", filepath.Base(path))
	fmt.Println(string(content))
}
