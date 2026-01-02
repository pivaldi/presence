// Package main demonstrates how to configure gorm.io/gen to generate
// database models using the presence package for nullable fields.
//
// This example:
//   - Starts a PostgreSQL container using testcontainers-go
//   - Creates sample tables with various column types
//   - Generates Go models with presence.Of[T] for nullable columns
//   - Demonstrates the generated model structure
//
// Usage:
//
//	cd examples/gorm-gen
//	go run main.go
//
// Prerequisites:
//   - Docker must be running
package main

import (
	"log"
	"os"
)

func main() {
	outDir := generate() // See file generate.go
	defer os.RemoveAll(outDir)

	display(outDir)

	printBar()
	log.Println("Example complete!")
}
