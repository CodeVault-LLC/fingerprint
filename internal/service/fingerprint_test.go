package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/codevault-llc/fingerprint/internal/service/models/entities"
	"github.com/codevault-llc/fingerprint/internal/service/models/repository"
	"github.com/codevault-llc/fingerprint/pkg/logger"
	pb "github.com/codevault-llc/fingerprint/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis      *bufconn.Listener
	mockRepo *repository.MockFingerprintRepo
)

func init() {
	lis = bufconn.Listen(bufSize)
	_, err := logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}

	go func() {
		s := grpc.NewServer()
		mockRepo = new(repository.MockFingerprintRepo)
		pb.RegisterFingerprintServiceServer(s, MockNewFingerprintService())
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAddFingerprint(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	mockRepo.On("AddFingerprint", mock.Anything).Return(&entities.Fingerprint{
		Id:          "123",
		Name:        "test",
		Description: "test",
		Pattern:     "testdata",
		Type:        "SCRIPT",
		Keywords:    []string{"test"},
	}, nil)

	client := pb.NewFingerprintServiceClient(conn)
	req := &pb.AddFingerprintRequest{
		Name:        "test",
		Description: "test",
		Pattern:     "testdata",
		Type:        pb.FingerprintType_SCRIPT,
		Keywords:    []string{"test"},
	}
	resp, err := client.AddFingerprint(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestGetFingerprint(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewFingerprintServiceClient(conn)
	addReq := &pb.AddFingerprintRequest{
		Name:        "test",
		Description: "test",
		Pattern:     "testdata",
		Type:        pb.FingerprintType_SCRIPT,
		Keywords:    []string{"test"},
	}
	_, err = client.AddFingerprint(ctx, addReq)
	assert.NoError(t, err)

	getReq := &pb.GetFingerprintRequest{Id: "123"}
	resp, err := client.GetFingerprint(ctx, getReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "testdata", resp.Pattern)
}
