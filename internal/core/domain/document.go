package domain

import (
	"time"

	proto "github.com/go-edi-document-processor/api/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Document struct {
	DocId      string
	Type       DocumentType
	Content    []byte
	SenderID   string
	ReceiverID string
	Status     DocumentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DocumentType string

const (
	XML  DocumentType = "XML"
	PDF  DocumentType = "PDF"
	JSON DocumentType = "JSON"
)

type DocumentStatus string

const (
	PENDING    DocumentStatus = "PENDING"
	RECEIVED   DocumentStatus = "RECEIVED"
	PROCESSED  DocumentStatus = "PROCESSED"
	FAILED     DocumentStatus = "FAILED"
	SUCCESSFUL DocumentStatus = "SUCCESSFUL"
)

func ProtoDocumentTypeToDomain(pt proto.DocumentType) DocumentType {
	switch pt {
	case proto.DocumentType_DOC_TYPE_XML:
		return XML
	case proto.DocumentType_DOC_TYPE_PDF:
		return PDF
	case proto.DocumentType_DOC_TYPE_JSON:
		return JSON
	default:
		return ""
	}
}

func DomainDocumentTypeToProto(dt DocumentType) proto.DocumentType {
	switch dt {
	case XML:
		return proto.DocumentType_DOC_TYPE_XML
	case PDF:
		return proto.DocumentType_DOC_TYPE_PDF
	case JSON:
		return proto.DocumentType_DOC_TYPE_JSON
	default:
		return proto.DocumentType_DOC_TYPE_UNSPECIFIED
	}
}

func ProtoDocumentStatusToDomain(ps proto.DocumentStatus) DocumentStatus {
	switch ps {
	case proto.DocumentStatus_DOC_STATUS_PENDING:
		return PENDING
	case proto.DocumentStatus_DOC_STATUS_RECEIVED:
		return RECEIVED
	case proto.DocumentStatus_DOC_STATUS_PROCESSED:
		return PROCESSED
	case proto.DocumentStatus_DOC_STATUS_FAILED:
		return FAILED
	case proto.DocumentStatus_DOC_STATUS_SUCCESSFUL:
		return SUCCESSFUL
	default:
		return ""
	}
}

func TimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func ProtoToDomainDocument(pdoc *proto.Document) *Document {
	if pdoc == nil {
		return nil
	}

	return &Document{
		DocId:      pdoc.DocId,
		Type:       ProtoDocumentTypeToDomain(pdoc.Type),
		Content:    pdoc.Content,
		SenderID:   pdoc.SenderId,
		ReceiverID: pdoc.ReceiverId,
		Status:     ProtoDocumentStatusToDomain(pdoc.Status),
		CreatedAt:  TimestampToTime(pdoc.CreatedAt),
		UpdatedAt:  TimestampToTime(pdoc.UpdatedAt),
	}
}

func DomainToProtoDocument(doc *Document) *proto.Document {
	if doc == nil {
		return nil
	}

	return &proto.Document{
		DocId:      doc.DocId,
		Type:       DomainDocumentTypeToProto(doc.Type),
		Content:    doc.Content,
		SenderId:   doc.SenderID,
		ReceiverId: doc.ReceiverID,
		Status:     DomainDocumentStatusToProto(doc.Status),
		CreatedAt:  TimeToTimestamp(doc.CreatedAt),
		UpdatedAt:  TimeToTimestamp(doc.UpdatedAt),
	}
}

func DomainDocumentStatusToProto(ds DocumentStatus) proto.DocumentStatus {
	switch ds {
	case PENDING:
		return proto.DocumentStatus_DOC_STATUS_PENDING
	case RECEIVED:
		return proto.DocumentStatus_DOC_STATUS_RECEIVED
	case PROCESSED:
		return proto.DocumentStatus_DOC_STATUS_PROCESSED
	case FAILED:
		return proto.DocumentStatus_DOC_STATUS_FAILED
	case SUCCESSFUL:
		return proto.DocumentStatus_DOC_STATUS_SUCCESSFUL
	default:
		return proto.DocumentStatus_DOC_STATUS_UNSPECIFIED
	}
}
