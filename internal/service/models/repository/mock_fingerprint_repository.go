package repository

import (
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/stretchr/testify/mock"
)

type MockFingerprintRepo struct {
	mock.Mock
}

func MockNewfingerprintRepo() FingerprintRepository {
	return &MockFingerprintRepo{}
}

func (m *MockFingerprintRepo) AddFingerprint(fingerprint *entities.Fingerprint) (*entities.Fingerprint, error) {
	args := m.Called(fingerprint)
	return args.Get(0).(*entities.Fingerprint), args.Error(1)
}

func (m *MockFingerprintRepo) GetFingerprint(id string) (*entities.Fingerprint, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Fingerprint), args.Error(1)
}

func (m *MockFingerprintRepo) MatchFingerprint(source string) ([]*entities.Fingerprint, error) {
	args := m.Called(source)
	return args.Get(0).([]*entities.Fingerprint), args.Error(1)
}
