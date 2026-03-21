PROTO_DIR := .
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT_DIR := .
GW_OUT_DIR := .

PROTOC := protoc
PROTOC_GO_PLUGIN := protoc-gen-go
PROTOC_GO_GRPC_PLUGIN := protoc-gen-go-grpc
PROTOC_GW_PLUGIN := protoc-gen-grpc-gateway

GO_OPT := paths=source_relative
GW_OPT := paths=source_relative
INCLUDE_PATHS := -I. -I$(GOPATH)/src/github.com/googleapis/googleapis

.PHONY: all clean

all: generate

generate: $(PROTO_FILES)
	$(PROTOC) $(INCLUDE_PATHS) \
		--go_out=$(GO_OUT_DIR) --go_opt=$(GO_OPT) \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=$(GO_OPT) \
		--grpc-gateway_out=$(GW_OUT_DIR) --grpc-gateway_opt=$(GW_OPT) \
		$(PROTO_FILES)

clean:
	rm -f $(GO_OUT_DIR)/*.pb.go
	rm -f $(GO_OUT_DIR)/*_grpc.pb.go
	rm -f $(GW_OUT_DIR)/*.pb.gw.go
