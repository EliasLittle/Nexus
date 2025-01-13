# Makefile

# Variables
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GRPC_GO = protoc-gen-go-grpc
PROTO_FILES = proto/nexus.proto
OUTPUT_DIR = pkg

# Targets
.PHONY: all clean proto build-client build-server run-client run-server


# Clean up generated files
clean:
	rm -f $(OUTPUT_DIR)/*.pb.go

# Regenerate Go code from proto files
proto:
	@echo "Regenerating Go code from proto files..."
	$(PROTOC) --go_out=$(OUTPUT_DIR) --go_opt=paths=source_relative \
    --go-grpc_out=$(OUTPUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_FILES)

# Build the Go application based on the provided target
build:
	@echo "Building the Go application: $(target)"
	go build -o $(target) ./cmd/$(target)

# Define specific targets for client and server
client:
	$(MAKE) build target=nexus-client

server:
	$(MAKE) build target=nexus-server

yukon:
	rm -f yukon; $(MAKE) build target=yukon

all: client server yukon

# Run the client application
run-client: build-client
	@echo "Running the client application..."
	./nexus-client

# Run the server application
run-server: build-server
	@echo "Running the server application..."
	./nexus-server
