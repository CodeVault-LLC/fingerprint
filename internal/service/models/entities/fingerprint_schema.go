package entities

func CreateFingerprintSchema() string {
	schema := `CREATE TABLE IF NOT EXISTS fingerprints (
			id UUID DEFAULT generateUUIDv4(),
			name String,
			description String,
			pattern String,
			type String,
			keywords Array(String),
			created_at DateTime DEFAULT now(),
			updated_at DateTime DEFAULT now()
		) ENGINE = MergeTree()
		ORDER BY id;`

	return schema
}
