package updater

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/updater/models/repository"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"github.com/codevault-llc/fingerprint/pkg/types"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

type Updater struct {
	database   *database.Database
	repository repository.UpdateRepository
}

func NewUpdater(db *database.Database) *Updater {
	repository := repository.NewUpdateRepository(db)

	if !db.TableExists("fingerprints") {
		fingerprints, err := loadFingerprints()
		if err == nil {
			if err := repository.BulkUpdateFingerprints(fingerprints); err != nil {
				logger.Log.Error("Failed to update fingerprints", zap.Error(err))
			}
		}
	}

	return &Updater{database: db, repository: repository}
}

type FingerprintData struct {
	Entries []types.Fingerprint `json:"entries"`
}

func loadFingerprints() ([]types.Fingerprint, error) {
	absPath, err := os.Getwd()
	if err != nil {
		logger.Log.Error("Error getting working directory", zap.Error(err))
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Validate JSON schema
	if err := validateJSONSchema(absPath+"/data/schema.json", absPath+"/data/fingerprints.json"); err != nil {
		logger.Log.Error("JSON schema validation failed", zap.Error(err))
		return nil, err
	}

	// Read and unmarshal fingerprints
	data, err := loadFingerprintFile(absPath + "/data/fingerprints.json")
	if err != nil {
		logger.Log.Error("Failed to load fingerprints", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Fingerprints loaded successfully", zap.Int("count", len(data.Entries)))
	return data.Entries, nil
}

func validateJSONSchema(schemaPath, dataPath string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	dataLoader := gojsonschema.NewReferenceLoader("file://" + dataPath)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return fmt.Errorf("JSON schema validation error: %w", err)
	}
	if !result.Valid() {
		for _, desc := range result.Errors() {
			logger.Log.Error("Schema validation error", zap.String("description", desc.Description()))
		}
		return fmt.Errorf("JSON schema validation failed")
	}
	return nil
}

func loadFingerprintFile(filePath string) (FingerprintData, error) {
	var data FingerprintData

	file, err := os.Open(filePath)
	if err != nil {
		return data, fmt.Errorf("failed to open fingerprints file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to unmarshal fingerprints JSON: %w", err)
	}
	return data, nil
}
