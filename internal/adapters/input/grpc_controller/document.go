package grpc_controller

import (
	"context"

	proto "github.com/go-edi-document-processor/api/proto/gen"
	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProtoDocumentServiceServer struct {
	proto.UnimplementedDocumentServiceServer
	documentService ports.DocumentService
}

func NewProtoDocumentServiceServer(documentService ports.DocumentService) proto.DocumentServiceServer {
	return &ProtoDocumentServiceServer{
		documentService: documentService,
	}
}

func (s *ProtoDocumentServiceServer) SendDocument(ctx context.Context, req *proto.SendDocumentRequest) (*proto.SendDocumentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domainDoc := domain.ProtoToDomainDocument(req.Document)

	savedDoc, err := s.documentService.SendDocument(ctx, domainDoc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.SendDocumentResponse{
		DocId:  savedDoc.DocId,
		Status: domain.DomainDocumentStatusToProto(savedDoc.Status),
	}, nil
}

func (s *ProtoDocumentServiceServer) GetDocumentByUUID(ctx context.Context, req *proto.GetDocumentByUUIDRequest) (*proto.GetDocumentByUUIDResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domainDoc, err := s.documentService.GetDocumentByUUID(ctx, req.DocId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if domainDoc == nil {
		return nil, status.Error(codes.NotFound, "document not found")
	}

	protoDoc := domain.DomainToProtoDocument(domainDoc)
	return &proto.GetDocumentByUUIDResponse{
		Document: protoDoc,
	}, nil
}

func (s *ProtoDocumentServiceServer) ReceiveDocument(ctx context.Context, req *proto.ReceiveDocumentRequest) (*proto.ReceiveDocumentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domainDoc, err := s.documentService.ReceiveDocument(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if domainDoc == nil {
		return nil, status.Error(codes.NotFound, "no document available for receiving")
	}

	protoDoc := domain.DomainToProtoDocument(domainDoc)
	return &proto.ReceiveDocumentResponse{
		Document: protoDoc,
	}, nil
}
