package service

import "github.com/codevault-llc/fingerprint/internal/service/models/repository"

type MockFingerprintRepo struct {
}

func MockNewFingerprintService() *FingerprintService {
	service := &FingerprintService{
		repository: repository.MockNewfingerprintRepo(),
	}

	return service
}
