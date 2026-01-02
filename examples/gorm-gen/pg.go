package main

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// SQL to create sample tables for demonstration
const initSQL = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    age INTEGER,
    balance NUMERIC(10, 2),
    is_active BOOLEAN NOT NULL DEFAULT true,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT false,
    view_count BIGINT DEFAULT 0,
    rating FLOAT,
    tags JSONB,
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id),
    bio TEXT,
    website VARCHAR(255),
    avatar_url VARCHAR(512),
    settings JSON,
    birth_date DATE,
    last_login TIMESTAMP
);

-- Insert sample data
INSERT INTO users (username, email, age, is_active, metadata) VALUES
    ('john_doe', 'john@example.com', 30, true, '{"role": "admin"}'),
    ('jane_smith', NULL, NULL, true, NULL);
`

func pgUp(ctx context.Context) (*postgres.PostgresContainer, *gorm.DB) {
	log.Println("Starting PostgreSQL container...")

	// Start PostgreSQL container using testcontainers
	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	log.Println("✓ PostgreSQL container started")

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	log.Printf("  Connection: %s", connStr)

	// Connect with GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(pg.Open(connStr), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create sample tables
	log.Println("Creating sample tables...")
	if err := db.Exec(initSQL).Error; err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	log.Println("✓ Sample tables created (users, posts, profiles)")

	return container, db
}
