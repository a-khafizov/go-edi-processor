package grpc

import (
	"context"

	"github.com/go-edi-document-processor/api/proto"
	"github.com/go-edi-document-processor/internal/infrastructure/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type documentService struct {
	proto.UnimplementedDocumentServiceServer
	log *logger.Logger
}

func NewDocumentService(log *logger.Logger) proto.DocumentServiceServer {
	return &documentService{log: log}
}

func (s *documentService) SendDocument(ctx context.Context, req *proto.SendDocumentRequest) (*proto.SendDocumentResponse, error) {
	s.log.Zap().Info("gRPC SendDocument called", zap.Any("request", req))
	return &proto.SendDocumentResponse{
		DocumentId: "generated-" + req.GetDocument().GetId(),
	}, nil
}

func (s *documentService) GetDocumentByUUID(ctx context.Context, req *proto.GetDocumentByUUIDRequest) (*proto.GetDocumentByUUIDResponse, error) {
	doc := &proto.Document{
		Id:         uuid.New().String(),
		Type:       proto.DocumentType_DOCUMENT_TYPE_XML,
		Content:    []byte("<xml>sample</xml>"),
		SenderId:   "sender-1",
		ReceiverId: "receiver-1",
		Metadata:   map[string]string{"key": "value"},
		Status:     proto.DocumentStatus_DOCUMENT_STATUS_PENDING,
	}
	return &proto.GetDocumentByUUIDResponse{
		Document: doc,
	}, nil
}
