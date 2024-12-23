package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/codevault-llc/fingerprint/config"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"go.uber.org/zap"
)

type Database struct {
	Db clickhouse.Conn
}

// validateDatabase checks if the ClickHouse database is accessible.
func validateDatabase() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Config.DatabaseAddr},
		Auth: clickhouse.Auth{
			Username: config.Config.DatabaseUser,
			Password: config.Config.DatabasePass,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	if err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.Config.DatabaseName)); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// NewDatabase initializes a new ClickHouse connection and ensures the schema exists.
func NewDatabase() (*Database, error) {
	validateDatabase()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.Config.DatabaseAddr},
		Auth: clickhouse.Auth{
			Database: config.Config.DatabaseName,
			Username: config.Config.DatabaseUser,
			Password: config.Config.DatabasePass,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "codevault-fingerprint", Version: "0.1"},
			},
		},
		// Disable TLS if not required
		TLS: nil,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
		Debugf:      log.Printf,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	db := &Database{Db: conn}

	if err := db.createSchema(); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

// BulkInsert performs a batch insert of fingerprints.
func (d *Database) BulkInsert(fingerprints []entities.Fingerprint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	batch, err := d.Db.PrepareBatch(ctx, entities.InsertFingerprintQuery())
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, fingerprint := range fingerprints {
		err := batch.Append(fingerprint.Id, fingerprint.Name, fingerprint.Description, fingerprint.Pattern, string(fingerprint.Type), fingerprint.Keywords, fingerprint.CreatedAt, fingerprint.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to append fingerprint data: %w", err)
		}
	}

	if err := batch.Flush(); err != nil {
		return fmt.Errorf("failed to flush batch: %w", err)
	}

	return nil
}

// createSchema creates the schema for the fingerprint table.
func (d *Database) createSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Db.Exec(ctx, entities.CreateFingerprintSchema()); err != nil {
		return fmt.Errorf("failed to execute schema creation query: %w", err)
	}

	return nil
}

// TableExists checks if a table exists in the database.
func (d *Database) TableExists(table string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count *uint64

	if err := d.Db.QueryRow(ctx, fmt.Sprintf("SELECT count() FROM %s", table)).Scan(&count); err != nil {
		logger.Log.Error("Failed to check if table exists", zap.Error(err))
		return false
	}

	return count != nil && *count > 0
}

// Close closes the database connection.
func (d *Database) Close() error {
	if d.Db != nil {
		if err := d.Db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}
	return nil
}
