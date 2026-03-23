.PHONY: all clean generate

PROTO_FILE = api/proto/document.proto
PROTO_PATHS = --proto_path=api/proto/ --proto_path=$(GOPATH)/pkg/mod/googleapis
OUTPUT_DIR = api/proto
OUTPUT_DIR_SWAGGER = docs
GO_OPTS = --go_opt=paths=source_relative
GRPC_OPTS = --go-grpc_opt=paths=source_relative
GATEWAY_OPTS = --grpc-gateway_opt=paths=source_relative
VALIDATE_OPTS = --validate_out="lang=go,paths=source_relative:$(OUTPUT_DIR)"

all: generate

generate: $(PROTO_FILE)
	protoc $(PROTO_PATHS) --go_out=$(OUTPUT_DIR) $(GO_OPTS) $(PROTO_FILE)
	protoc $(PROTO_PATHS) --go-grpc_out=$(OUTPUT_DIR) $(GRPC_OPTS) $(PROTO_FILE)
	protoc $(PROTO_PATHS) --grpc-gateway_out=$(OUTPUT_DIR) $(GATEWAY_OPTS) $(PROTO_FILE)
	protoc $(PROTO_PATHS) --openapiv2_out=$(OUTPUT_DIR_SWAGGER) $(PROTO_FILE)
	protoc $(PROTO_PATHS) $(VALIDATE_OPTS) $(PROTO_FILE)

clean:
	rm -f $(OUTPUT_DIR)/*.pb.go $(OUTPUT_DIR)/*.swagger.json $(OUTPUT_DIR)/*.pb.gw.go $(OUTPUT_DIR)/*_validate.pb.go
