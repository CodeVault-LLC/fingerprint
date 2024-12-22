package repository

import (
	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/google/uuid"
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

func (r *FingerprintRepo) AddFingerprint(fingerprint *entities.Fingerprint) (*entities.Fingerprint, error) {
	fingerprint.Id = uuid.New().String()

	query := `INSERT INTO fingerprint.fingerprints (id, name, description, pattern, type, keywords) VALUES (?, ?, ?, ?, ?, ?)`

	err := r.Db.Db.Query(query, fingerprint.Id, fingerprint.Name, fingerprint.Description, fingerprint.Pattern, string(fingerprint.Type), fingerprint.Keywords).Exec()
	if err != nil {
		return nil, err
	}

	return fingerprint, nil
}

func (r *FingerprintRepo) GetFingerprint(id string) (*entities.Fingerprint, error) {
	var fingerprint entities.Fingerprint
	query := `SELECT id, name, description, pattern, type, keywords, created_at, updated_at FROM fingerprint.fingerprints WHERE id = ?`
	err := r.Db.Db.Query(query, id).Scan(&fingerprint.Id, &fingerprint.Name, &fingerprint.Description, &fingerprint.Pattern, &fingerprint.Type, &fingerprint.Keywords, &fingerprint.CreatedAt, &fingerprint.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &fingerprint, nil
}

func (r *FingerprintRepo) MatchFingerprint(source string) ([]*entities.Fingerprint, error) {
	var fingerprints []*entities.Fingerprint
	query := `SELECT id, name, description, pattern, type, keywords, created_at, updated_at FROM fingerprint.fingerprints WHERE ? LIKE CONCAT('%', pattern, '%')`
	iter := r.Db.Db.Query(query, source).Iter()
	for {
		var fingerprint entities.Fingerprint
		if !iter.Scan(&fingerprint.Id, &fingerprint.Name, &fingerprint.Description, &fingerprint.Pattern, &fingerprint.Type, &fingerprint.Keywords, &fingerprint.CreatedAt, &fingerprint.UpdatedAt) {
			break
		}
		fingerprints = append(fingerprints, &fingerprint)
	}

	return fingerprints, nil
}
