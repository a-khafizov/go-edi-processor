package adapters

import (
	"context"
	"time"

	proto "github.com/go-edi-document-processor/api/proto/gen"
	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/primary"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	domainDoc := protoToDomainDocument(req.Document)

	savedDoc, err := s.documentService.SendDocument(ctx, domainDoc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.SendDocumentResponse{
		DocId:  savedDoc.DocId,
		Status: domainDocumentStatusToProto(savedDoc.Status),
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

	protoDoc := domainToProtoDocument(domainDoc)
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

	protoDoc := domainToProtoDocument(domainDoc)
	return &proto.ReceiveDocumentResponse{
		Document: protoDoc,
	}, nil
}

func protoDocumentTypeToDomain(pt proto.DocumentType) domain.DocumentType {
	switch pt {
	case proto.DocumentType_DOC_TYPE_XML:
		return domain.XML
	case proto.DocumentType_DOC_TYPE_PDF:
		return domain.PDF
	case proto.DocumentType_DOC_TYPE_JSON:
		return domain.JSON
	default:
		return ""
	}
}

func domainDocumentTypeToProto(dt domain.DocumentType) proto.DocumentType {
	switch dt {
	case domain.XML:
		return proto.DocumentType_DOC_TYPE_XML
	case domain.PDF:
		return proto.DocumentType_DOC_TYPE_PDF
	case domain.JSON:
		return proto.DocumentType_DOC_TYPE_JSON
	default:
		return proto.DocumentType_DOC_TYPE_UNSPECIFIED
	}
}

func protoDocumentStatusToDomain(ps proto.DocumentStatus) domain.DocumentStatus {
	switch ps {
	case proto.DocumentStatus_DOC_STATUS_PENDING:
		return domain.Pending
	case proto.DocumentStatus_DOC_STATUS_RECEIVED:
		return domain.Received
	case proto.DocumentStatus_DOC_STATUS_PROCESSED:
		return domain.Processed
	case proto.DocumentStatus_DOC_STATUS_FAILED:
		return domain.Failed
	case proto.DocumentStatus_DOC_STATUS_SUCCESSFUL:
		return domain.Successful
	default:
		return ""
	}
}

func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func protoToDomainDocument(pdoc *proto.Document) *domain.Document {
	if pdoc == nil {
		return nil
	}

	return &domain.Document{
		DocId:      pdoc.DocId,
		Type:       protoDocumentTypeToDomain(pdoc.Type),
		Content:    pdoc.Content,
		SenderID:   pdoc.SenderId,
		ReceiverID: pdoc.ReceiverId,
		Status:     protoDocumentStatusToDomain(pdoc.Status),
		CreatedAt:  timestampToTime(pdoc.CreatedAt),
		UpdatedAt:  timestampToTime(pdoc.UpdatedAt),
	}
}

func domainToProtoDocument(doc *domain.Document) *proto.Document {
	if doc == nil {
		return nil
	}

	return &proto.Document{
		DocId:      doc.DocId,
		Type:       domainDocumentTypeToProto(doc.Type),
		Content:    doc.Content,
		SenderId:   doc.SenderID,
		ReceiverId: doc.ReceiverID,
		Status:     domainDocumentStatusToProto(doc.Status),
		CreatedAt:  timeToTimestamp(doc.CreatedAt),
		UpdatedAt:  timeToTimestamp(doc.UpdatedAt),
	}
}

func domainDocumentStatusToProto(ds domain.DocumentStatus) proto.DocumentStatus {
	switch ds {
	case domain.Pending:
		return proto.DocumentStatus_DOC_STATUS_PENDING
	case domain.Received:
		return proto.DocumentStatus_DOC_STATUS_RECEIVED
	case domain.Processed:
		return proto.DocumentStatus_DOC_STATUS_PROCESSED
	case domain.Failed:
		return proto.DocumentStatus_DOC_STATUS_FAILED
	case domain.Successful:
		return proto.DocumentStatus_DOC_STATUS_SUCCESSFUL
	default:
		return proto.DocumentStatus_DOC_STATUS_UNSPECIFIED
	}
}
