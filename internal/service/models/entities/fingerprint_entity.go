package entities

import "time"

type FingerprintType string

const (
	Script FingerprintType = "script"
)

type Fingerprint struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Pattern     string          `json:"pattern"`
	Type        FingerprintType `json:"type"`
	Keywords    []string        `json:"keywords"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
