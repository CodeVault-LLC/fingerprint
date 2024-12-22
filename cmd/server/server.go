package server

import (
	"fmt"
	"net"

	"github.com/codevault-llc/fingerprint/config"
	"github.com/codevault-llc/fingerprint/internal/database"
	"github.com/codevault-llc/fingerprint/internal/service"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/codevault-llc/fingerprint/proto"
)

func StartServer(config *config.InternalConfig, database *database.Database) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.ServiceHost, config.ServicePort))
	if err != nil {
		logger.Log.Error("Error starting listener", zap.Error(err))
		return err
	}

	grpcServer := grpc.NewServer()
	fingerprintService := service.NewFingerprintService(database)

	pb.RegisterFingerprintServiceServer(grpcServer, fingerprintService)

	logger.Log.Info("Starting server", zap.String("host", config.ServiceHost), zap.String("port", config.ServicePort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Log.Error("Error starting server", zap.Error(err))
		return err
	}

	return nil
}
