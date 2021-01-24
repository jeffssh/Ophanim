package message

import (
	context "context"
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

/*
type Confidence int
type MessageType int

// Message and Confidence types
const (
	LOW Confidence = iota
	MEDIUM
	HIGH
	CERTAIN

	INFORMATIONAL MessageType = iota
	VULNERABILITY
)

func (c Confidence) String() string {
	return [...]string{"Low", "Medium", "High", "Certain"}[c]
}

func (t MessageType) String() string {
	return [...]string{"Informational", "Vulnerability"}[t]
}

func (m *Message) String() string {

	switch m.MessageType {
	case VULNERABILITY:
		return fmt.Sprintf("Vulnerability - Confidence: %s\nDescription: %s", m.Confidence.String(), m.Message)
	default:
		return m.String()
	}
}

// Message - Messages are sent to module named pipes by module processes
type Message struct {
	Confidence  Confidence
	MessageType MessageType
	Message     []string
}
*/
type reportMessageServer struct {
	ReportMessageServiceServer
}

func (s *reportMessageServer) SendReportMessage(ctx context.Context, in *ReportMessage) (*emptypb.Empty, error) {
	log.Printf("Received: %v", in)
	return &emptypb.Empty{}, nil
}

// NewMessageServer - create new gRPC server
func NewMessageServer() (s *grpc.Server) {
	//lis, err := net.Listen(network, address)
	s = grpc.NewServer()
	RegisterReportMessageServiceServer(s, &reportMessageServer{})
	return
}
