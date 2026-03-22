package grpc_controllers

import (
	"context"
	"fmt"

	"github.com/go-edi-document-processor/api/proto"
	"github.com/go-edi-document-processor/internal/controllers/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type documentService struct {
	proto.UnimplementedDocumentServiceServer
	logger *zap.Logger
}

func NewDocumentService(logger *zap.Logger) proto.DocumentServiceServer {
	return &documentService{logger: logger}
}

func (s *documentService) SendDocument(ctx context.Context, req *proto.SendDocumentRequest) (*proto.SendDocumentResponse, error) {
	if req.Document == nil {
		return nil, fmt.Errorf("document is required %s", req.Document.Type)
	}

	doc := &dto.Document{
		Type:       req.Document.Type,
		Content:    req.Document.Content,
		SenderID:   req.Document.SenderId,
		ReceiverID: req.Document.ReceiverId,
		Status:     req.Document.Status,
	}

	return &proto.SendDocumentResponse{
		DocumentId: doc.DocID,
	}, nil
}

func (s *documentService) GetDocumentByUUID(ctx context.Context, req *proto.GetDocumentByUUIDRequest) (*proto.GetDocumentByUUIDResponse, error) {
	doc := &proto.Document{
		DocId:      uuid.New().String(),
		Type:       proto.DocumentType_DOCUMENT_TYPE_XML,
		Content:    []byte("<xml>sample</xml>"),
		SenderId:   "sender-1",
		ReceiverId: "receiver-1",
		Status:     proto.DocumentStatus_DOCUMENT_STATUS_PENDING,
	}
	return &proto.GetDocumentByUUIDResponse{
		Document: doc,
	}, nil
}
