package main

import (
	"fmt"
	"path/filepath"

	"gorm.io/gen"
	"gorm.io/gorm"
)

// dataTypeMap defines custom type mappings for PostgreSQL column types.
// For nullable columns, we wrap the type with presence.Of[T].
var dataTypeMap = map[string]func(gorm.ColumnType) string{
	// String types
	"varchar": stringMapFunc,
	"text":    stringMapFunc,
	"char":    stringMapFunc,
	"bpchar":  stringMapFunc,

	// Integer types
	"int2":     integerMapFunc,
	"int4":     integerMapFunc,
	"int8":     integerMapFunc,
	"smallint": integerMapFunc,
	"integer":  integerMapFunc,
	"bigint":   integerMapFunc,

	// Floating point types
	"float4":  floatMapFunc,
	"float8":  floatMapFunc,
	"real":    floatMapFunc,
	"numeric": floatMapFunc,
	"decimal": floatMapFunc,

	// Boolean
	"bool":    boolMapFunc,
	"boolean": boolMapFunc,

	// Date/Time types
	"date":        dateMapFunc,
	"time":        timeMapFunc,
	"timetz":      timeMapFunc,
	"timestamp":   timestampMapFunc,
	"timestamptz": timestampMapFunc,

	// JSON types
	"json":  jsonMapFunc,
	"jsonb": jsonMapFunc,

	// UUID
	"uuid": uuidMapFunc,
}

// Type mapping functions

func wrapNullable(c gorm.ColumnType, baseType string) string {
	if nullable, _ := c.Nullable(); nullable {
		return fmt.Sprintf("presence.Of[%s]", baseType)
	}
	return baseType
}

func stringMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "string")
}

func integerMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "int64")
}

func floatMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "float64")
}

func boolMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "bool")
}

func dateMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "time.Time")
}

func timeMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "time.Time")
}

func timestampMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "time.Time")
}

func jsonMapFunc(c gorm.ColumnType) string {
	// JSON columns always use presence.Of[any] for convenient handling
	return "presence.Of[any]"
}

func uuidMapFunc(c gorm.ColumnType) string {
	return wrapNullable(c, "uuid.UUID")
}

func getGenerator(outputDir string, db *gorm.DB) *gen.Generator {
	dalPath := filepath.Join(outputDir, "dal")
	modelPath := getModelOutDir(outputDir)

	// Generator configuration
	config := gen.Config{
		OutPath:      dalPath,
		ModelPkgPath: modelPath,
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,

		FieldCoverable:   true,
		FieldWithTypeTag: true,
		FieldNullable:    false, // We handle nullable via WithDataTypeMap
	}

	// Add required import paths for the generated code
	config.WithImportPkgPath(
		"github.com/pivaldi/presence",
		"github.com/google/uuid",
	)

	// Configure JSON tag naming strategy
	config.WithJSONTagNameStrategy(snakeToCamelCase)

	// Create the generator
	g := gen.NewGenerator(config)

	// Use the database connection
	g.UseDB(db)

	// Apply custom data type mappings
	g.WithDataTypeMap(dataTypeMap)

	// Config to generate models for all tables
	g.ApplyBasic(g.GenerateAllTable()...)

	return g
}
