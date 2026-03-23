package adapters

import (
	"context"
	"time"

	"github.com/go-edi-document-processor/api/proto"
	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/primary"
	"github.com/google/uuid"
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

	doc := &proto.Document{
		DocId:      uuid.New().String(),
		Type:       proto.DocumentType_DOC_TYPE_XML,
		Content:    []byte("<xml>sample</xml>"),
		SenderId:   "sender-1",
		ReceiverId: "receiver-1",
		Status:     proto.DocumentStatus_DOC_STATUS_PENDING,
	}
	return &proto.GetDocumentByUUIDResponse{
		Document: doc,
	}, nil
}

func (s *ProtoDocumentServiceServer) ReceiveDocument(ctx context.Context, req *proto.ReceiveDocumentRequest) (*proto.ReceiveDocumentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc := &proto.Document{
		DocId:      uuid.New().String(),
		Type:       proto.DocumentType_DOC_TYPE_XML,
		Content:    []byte("<xml>sample</xml>"),
		SenderId:   "sender-1",
		ReceiverId: "receiver-1",
		Status:     proto.DocumentStatus_DOC_STATUS_PENDING,
	}
	return &proto.ReceiveDocumentResponse{
		Document: doc,
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

func protoToDomainDocument(pdoc *proto.Document) *domain.Document {
	if pdoc == nil {
		return nil
	}

	return &domain.Document{
		DocId:      pdoc.DocId,
		Type:       protoDocumentTypeToDomain(pdoc.Type),
		Content:    pdoc.Content, // bytes → []byte (совместимы)
		SenderID:   pdoc.SenderId,
		ReceiverID: pdoc.ReceiverId,
		Status:     protoDocumentStatusToDomain(pdoc.Status),
		CreatedAt:  timestampToTime(pdoc.CreatedAt),
		UpdatedAt:  timestampToTime(pdoc.UpdatedAt),
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
