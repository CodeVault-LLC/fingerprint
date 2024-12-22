package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/internal/service/models/repository"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	pb "github.com/codevault-llc/fingerprint/proto"
	"go.uber.org/zap"
)

// FingerprintService implements the FingerprintService gRPC server.
type FingerprintService struct {
	pb.UnimplementedFingerprintServiceServer
	mu         sync.RWMutex // RWMutex for better performance on reads
	repository repository.FingerprintRepository
}

// NewFingerprintService initializes the FingerprintService with default data.
func NewFingerprintService(database *database.Database) *FingerprintService {
	service := &FingerprintService{
		repository: repository.NewFingerprintRepository(database),
	}

	return service
}

// AddFingerprint adds a fingerprint to the server.
func (s *FingerprintService) AddFingerprint(ctx context.Context, req *pb.AddFingerprintRequest) (*pb.AddFingerprintResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fingerprint, err := s.repository.AddFingerprint(&entities.Fingerprint{
		Name:        req.Name,
		Description: req.Description,
		Pattern:     req.Pattern,
		Type:        entities.FingerprintType(req.Type),
		Keywords:    req.Keywords,
	})

	if err != nil {
		logger.Log.Error("Error adding fingerprint", zap.Error(err))
		return nil, err
	}

	return &pb.AddFingerprintResponse{Id: fingerprint.Id}, nil
}

// GetFingerprint retrieves a fingerprint by ID.
func (s *FingerprintService) GetFingerprint(ctx context.Context, req *pb.GetFingerprintRequest) (*pb.GetFingerprintResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fingerprint, err := s.repository.GetFingerprint(req.Id)
	if err != nil {
		logger.Log.Error("Error getting fingerprint", zap.Error(err))
		return nil, err
	}

	return &pb.GetFingerprintResponse{
		Id:          fingerprint.Id,
		Name:        fingerprint.Name,
		Description: fingerprint.Description,
		Pattern:     fingerprint.Pattern,
		Type:        string(fingerprint.Type),
		Keywords:    fingerprint.Keywords,
		CreatedAt:   fingerprint.CreatedAt.String(),
		UpdatedAt:   fingerprint.UpdatedAt.String(),
	}, nil
}

// MatchFingerprint matches a script against stored fingerprints.
func (s *FingerprintService) MatchFingerprint(ctx context.Context, req *pb.MatchFingerprintRequest) (*pb.MatchFingerprintResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return nil, fmt.Errorf("no match found")
}
