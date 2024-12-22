package repository

import (
	"time"

	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/pkg/types"
	"github.com/google/uuid"
)

type UpdateRepository interface {
	BulkUpdateFingerprints(fingerprints []types.Fingerprint) error
}

type UpdateRepo struct {
	Db *database.Database
}

func NewUpdateRepository(db *database.Database) UpdateRepository {
	return &UpdateRepo{Db: db}
}

func (r *UpdateRepo) BulkUpdateFingerprints(fingerprints []types.Fingerprint) error {
	bulkFingerprints := make([]entities.Fingerprint, 0, len(fingerprints))
	for _, fingerprint := range fingerprints {
		bulkFingerprints = append(bulkFingerprints, entities.Fingerprint{
			Id:          uuid.New().String(),
			Name:        fingerprint.Name,
			Description: fingerprint.Description,
			Pattern:     fingerprint.Regex,
			Type:        entities.FingerprintType(fingerprint.Type),
			Keywords:    fingerprint.Keywords,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	if err := r.Db.BulkInsert(bulkFingerprints); err != nil {
		return err
	}

	return nil
}
