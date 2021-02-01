package message

import (
	context "context"
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reportMessageServer struct {
	ReportMessageServiceServer
}

func (s *reportMessageServer) SendReportMessage(ctx context.Context, in *ReportMessage) (*emptypb.Empty, error) {
	log.Printf("Received: %v", in)
	return &emptypb.Empty{}, nil
}

// NewMessageServer - create new gRPC server
func NewMessageServer() (s *grpc.Server) {
	s = grpc.NewServer()
	RegisterReportMessageServiceServer(s, &reportMessageServer{})
	return
}
