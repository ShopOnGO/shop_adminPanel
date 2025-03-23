PROTO_DIR=internal/proto
OUT_DIR=pkg

export PATH := $(PATH):$(shell go env GOPATH)/bin

generate:
	@mkdir -p $(OUT_DIR)
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/*.proto

generate_brand:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/brand.proto

generate_home:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/home.proto

generate_entities:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/entities.proto

generate_links:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/link.proto

generate_categories:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/category.proto

generate_products:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/product.proto

generate_users:
	@protoc --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_DIR)/user.proto
