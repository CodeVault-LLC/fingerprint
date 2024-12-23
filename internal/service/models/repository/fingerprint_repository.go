package repository

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FingerprintRepository interface {
	AddFingerprint(fingerprint *entities.Fingerprint) (*entities.Fingerprint, error)
	GetFingerprint(id string) (*entities.Fingerprint, error)
	MatchFingerprint(source string) ([]*entities.Fingerprint, error)
}

type FingerprintRepo struct {
	Db *database.Database
}

func NewFingerprintRepository(db *database.Database) FingerprintRepository {
	return &FingerprintRepo{Db: db}
}

// AddFingerprint inserts a new fingerprint into the database.
func (r *FingerprintRepo) AddFingerprint(fingerprint *entities.Fingerprint) (*entities.Fingerprint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fingerprint.Id = uuid.New().String()
	query := `
		INSERT INTO fingerprint.fingerprints (id, name, description, pattern, type, keywords, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	err := r.Db.Db.Exec(ctx, query,
		fingerprint.Id,
		fingerprint.Name,
		fingerprint.Description,
		fingerprint.Pattern,
		string(fingerprint.Type),
		fingerprint.Keywords,
		fingerprint.CreatedAt,
		fingerprint.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add fingerprint: %w", err)
	}

	return fingerprint, nil
}

// GetFingerprint retrieves a fingerprint by ID.
func (r *FingerprintRepo) GetFingerprint(id string) (*entities.Fingerprint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT id, name, description, pattern, type, keywords, created_at, updated_at
		FROM fingerprint.fingerprints
		WHERE id = ?
	`

	var fingerprint entities.Fingerprint
	err := r.Db.Db.QueryRow(ctx, query, id).Scan(
		&fingerprint.Id,
		&fingerprint.Name,
		&fingerprint.Description,
		&fingerprint.Pattern,
		&fingerprint.Type,
		&fingerprint.Keywords,
		&fingerprint.CreatedAt,
		&fingerprint.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fingerprint: %w", err)
	}

	return &fingerprint, nil
}

// MatchFingerprint searches for fingerprints that match the source string.
func (r *FingerprintRepo) MatchFingerprint(source string) ([]*entities.Fingerprint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT id, name, description, pattern, type, keywords, created_at, updated_at
		FROM fingerprint.fingerprints
	`

	rows, err := r.Db.Db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query fingerprints: %w", err)
	}
	defer rows.Close()

	var fingerprints []*entities.Fingerprint
	for rows.Next() {
		var fingerprint entities.SafeFingerprint
		err := rows.Scan(
			&fingerprint.Id,
			&fingerprint.Name,
			&fingerprint.Description,
			&fingerprint.Pattern,
			&fingerprint.Type,
			&fingerprint.Keywords,
			&fingerprint.CreatedAt,
			&fingerprint.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fingerprint: %w", err)
		}

		// Perform regex matching in Go
		matched, matchErr := regexp.MatchString(fingerprint.Pattern, source)
		if matchErr != nil {
			logger.Log.Error("Failed to evaluate regex", zap.String("pattern", fingerprint.Pattern), zap.Error(matchErr))
			continue
		}

		if matched {
			fingerprints = append(fingerprints, &entities.Fingerprint{
				Id:          fingerprint.Id,
				Name:        fingerprint.Name,
				Description: fingerprint.Description,
				Pattern:     fingerprint.Pattern,
				Type:        entities.FingerprintType(fingerprint.Type),
				Keywords:    fingerprint.Keywords,
				CreatedAt:   fingerprint.CreatedAt,
				UpdatedAt:   fingerprint.UpdatedAt,
			})
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error during row iteration: %w", rows.Err())
	}

	return fingerprints, nil
}
