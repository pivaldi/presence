package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
)

func generate() string {
	ctx := context.Background()
	container, db := pgUp(ctx)
	defer func() {
		log.Println("Terminating PostgreSQL container...")
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}()

	// Output paths for generated code
	outDir, err := os.MkdirTemp("", "gorm-gen-example-*")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	log.Println("Generating models…")
	getGenerator(outDir, db).Execute()

	log.Println("✓ Models generated successfully!")
	log.Printf("  Output directory: %s", outDir)

	return outDir
}

func getModelOutDir(outDir string) string {
	return filepath.Join(outDir, "models")
}
