package seeds

import "sync"

type Seed struct {
	mu           sync.RWMutex
	fingerprints map[string]string
}

func NewSeed() *Seed {
	seed := &Seed{
		fingerprints: make(map[string]string),
	}

	seed.initializeFingerprints()
	return seed
}

func (s *Seed) initializeFingerprints() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.fingerprints["1"] = "pattern1"
	s.fingerprints["2"] = "pattern2"
	s.fingerprints["3"] = "pattern3"
}
