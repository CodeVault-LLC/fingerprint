package entities

import "github.com/codevault-llc/fingerprint/config"

func CreateFingerprintSchema() string {
	keyspace := config.Config.DatabaseKeyspace

	schema := `CREATE TABLE IF NOT EXISTS ` + keyspace + `.fingerprints (
		id UUID PRIMARY KEY,
		name TEXT,
		description TEXT,
		pattern TEXT,
		type TEXT,
		keywords SET<TEXT>,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);`

	return schema
}
