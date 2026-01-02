// Package main demonstrates how to configure gorm.io/gen to generate
// database models using the presence package for presence fields.
//
// This example shows:
//   - How to configure gen.Config for presence field generation
//   - Custom type mapping for PostgreSQL types (json, jsonb, uuid, date)
//   - Using WithPresenceNameStrategy to wrap presence fields with presence.Of[T]
//
// Usage:
//
//	go run main.go
//
// Prerequisites:
//   - A running PostgreSQL database
//   - Update the DSN connection string below
package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// dataTypeMap defines custom type mappings for PostgreSQL column types.
// These mappings control how database types are converted to Go types.
var dataTypeMap = map[string]func(gorm.ColumnType) string{
	"json":  jsonMapFunc,
	"jsonb": jsonMapFunc,
	"date":  dateMapFunc,
	"int2":  integerMapFunc,
	"int4":  integerMapFunc,
	"int8":  integerMapFunc,
	"int16": integerMapFunc,
	"int32": integerMapFunc,
	"uuid":  uuidMapFunc,
}

func main() {
	// Database connection string - update with your credentials
	dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"

	// Schema to generate models from
	schemaName := "public"

	// Output paths for generated code
	dalPath := "./generated/dal"
	modelPath := "./generated/models"

	// Generator configuration
	config := gen.Config{
		OutPath:      dalPath,   // Output path for query/DAL code
		ModelPkgPath: modelPath, // Output path for model structs
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
		OutFile:      schemaName + ".go",

		// Enable presence field generation
		FieldCoverable:   true, // Generate pointer fields for presence columns
		FieldWithTypeTag: true, // Add type tag to struct fields
		FieldPresence:    true, // Enable presence field generation
	}

	// Configure presence type wrapper strategy.
	// This wraps presence database fields with presence.Of[T].
	// For example, a presence VARCHAR becomes presence.Of[string]
	config.WithPresenceNameStrategy(func(fieldType string) string {
		return fmt.Sprintf("presence.Of[%s]", fieldType)
	})

	// Add required import paths for the generated code
	config.WithImportPkgPath(
		"github.com/pivaldi/presence",
		"github.com/google/uuid",
	)

	// Configure JSON tag naming strategy (snake_case to camelCase)
	config.WithJSONTagNameStrategy(snakeToCamelCase)

	// Create the generator
	g := gen.NewGenerator(config)

	// GORM configuration with schema naming strategy
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schemaName + ".", // PostgreSQL schema prefix
			SingularTable: true,             // Use singular table names (user instead of users)
		},
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Use the database connection for model generation
	g.UseDB(db)

	// Apply custom data type mappings for PostgreSQL types
	g.WithDataTypeMap(dataTypeMap)

	// Generate models for all tables in the schema
	g.ApplyBasic(g.GenerateAllTable()...)

	// Alternative: Generate specific tables only
	// g.ApplyBasic(
	//     g.GenerateModel("users"),
	//     g.GenerateModel("posts"),
	//     g.GenerateModelAs("user_profiles", "Profile"), // Custom model name
	// )

	// Alternative: Generate with custom field options
	// g.ApplyBasic(
	//     g.GenerateModel("users",
	//         gen.FieldIgnore("password_hash"),           // Ignore specific field
	//         gen.FieldType("status", "UserStatus"),      // Custom type for field
	//         gen.FieldGORMTag("email", "uniqueIndex"),   // Add GORM tag
	//     ),
	// )

	// Execute the generation
	g.Execute()

	fmt.Println("Models generated successfully!")
	fmt.Printf("  DAL code: %s\n", dalPath)
	fmt.Printf("  Models:   %s\n", modelPath)
}

// Type mapping functions for PostgreSQL types

// jsonMapFunc maps json/jsonb columns to appropriate Go types.
// Presence JSON columns will be wrapped by WithPresenceNameStrategy.
func jsonMapFunc(c gorm.ColumnType) string {
	if presence, _ := c.Presence(); presence {
		// Return base type - WithPresenceNameStrategy will wrap it as presence.Of[any]
		return "any"
	}
	// Non-presence JSON still uses presence.Of for convenient JSON handling
	return "presence.Of[any]"
}

// uuidMapFunc maps uuid columns to uuid.UUID type.
func uuidMapFunc(_ gorm.ColumnType) string {
	return "uuid.UUID"
}

// integerMapFunc maps all integer types to int64 for consistency.
func integerMapFunc(_ gorm.ColumnType) string {
	return "int64"
}

// dateMapFunc maps date/timestamp columns to time.Time.
// Presence date columns will be wrapped by WithPresenceNameStrategy.
func dateMapFunc(c gorm.ColumnType) string {
	if presence, _ := c.Presence(); presence {
		// Return base type - WithPresenceNameStrategy will wrap it as presence.Of[time.Time]
		return "time.Time"
	}
	// Non-presence dates still benefit from presence.Of for zero value handling
	return "presence.Of[time.Time]"
}

// snakeToCamelCase converts snake_case to camelCase for JSON tags.
// Example: "created_at" -> "createdAt"
func snakeToCamelCase(in string) string {
	in = strings.TrimSpace(cases.Lower(language.Und).String(in))
	if in == "" {
		return in
	}

	tokens := strings.Split(in, "_")
	caser := cases.Title(language.Und, cases.NoLower)

	var out string
	for i, token := range tokens {
		if i == 0 {
			out += token // First token stays lowercase
			continue
		}
		out += caser.String(token) // Capitalize subsequent tokens
	}

	return out
}
