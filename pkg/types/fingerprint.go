package types

type FingerprintType string

const (
	FingerprintTypeScript FingerprintType = "script"
)

type Fingerprint struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Regex       string          `json:"regex"`
	Type        FingerprintType `json:"type"`
	Keywords    []string        `json:"keywords"`
}
