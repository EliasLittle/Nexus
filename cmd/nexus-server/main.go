package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "nexus/pkg/proto"

	"google.golang.org/grpc"
)

const (
	port = 50051
)

// nexusServer implements the NexusService gRPC service
type nexusServer struct {
	pb.UnimplementedNexusServiceServer
}

// Initialize a new Trie for storing paths
var pathTrie = NewTrie()

// RegisterEventStream implements the publisher endpoint for registering event streams
func (s *nexusServer) RegisterEventStream(ctx context.Context, req *pb.RegisterEventStreamRequest) (*pb.RegisterEventStreamResponse, error) {
	log.Printf("Received event stream registration request for path: %s, topic: %s\n", req.Path, req.EventStream.Topic)

	// Insert the event stream into the Trie
	pathTrie.Insert(req.Path, req.EventStream) // Store the EventStream
	pathTrie.Traverse()                        // Print the Trie after the update

	return &pb.RegisterEventStreamResponse{Success: true}, nil
}

// RegisterDataset implements the publisher endpoint for registering datasets
func (s *nexusServer) RegisterDataset(ctx context.Context, req *pb.RegisterDatasetRequest) (*pb.RegisterDatasetResponse, error) {
	log.Printf("Received dataset registration request for path: %s\n", req.Path)

	// Insert the dataset into the Trie
	pathTrie.Insert(req.Path, req.Dataset) // Store the Dataset
	pathTrie.Traverse()                    // Print the Trie after the update

	return &pb.RegisterDatasetResponse{Success: true}, nil
}

// StoreValue implements the publisher endpoint for storing values
func (s *nexusServer) StoreValue(ctx context.Context, req *pb.StoreValueRequest) (*pb.StoreValueResponse, error) {
	log.Printf("Received value storage request for path: %s\n", req.Path)

	pathTrie.Insert(req.Path, req.DirectValue)

	//pathTrie.Traverse() // Print the Trie after the update
	return &pb.StoreValueResponse{Success: true}, nil
}

// GetValue implements the consumer endpoint for retrieving values
func (s *nexusServer) GetValue(ctx context.Context, req *pb.GetPathRequest) (*pb.GetValueResponse, error) {
	log.Printf("Received value request for path: %s\n", req.Path)

	value, err := pathTrie.Get(req.Path)

	if err != nil {
		return &pb.GetValueResponse{Error: err.Error()}, nil
	}

	// Check if the value is of type *pb.DirectValue
	directValue, ok := value.(*pb.DirectValue) // Use pointer type assertion
	if !ok {
		return &pb.GetValueResponse{Error: "value is not of type DirectValue"}, nil
	}

	return &pb.GetValueResponse{Value: directValue}, nil
}

// ListChildren implements the endpoint for listing children nodes
func (s *nexusServer) ListChildren(ctx context.Context, req *pb.ListChildrenRequest) (*pb.ListChildrenResponse, error) {
	log.Printf("Received request to list children for path: %s\n", req.Path)

	children := pathTrie.GetChildren(req.Path)

	return &pb.ListChildrenResponse{Children: children}, nil
}

// GetEventStream retrieves event stream details from the specified path
func (s *nexusServer) GetEventStream(ctx context.Context, req *pb.GetPathRequest) (*pb.GetEventStreamResponse, error) {
	log.Printf("Received event stream request for path: %s\n", req.Path)

	value, err := pathTrie.Get(req.Path)
	if err != nil {
		return &pb.GetEventStreamResponse{Error: err.Error()}, nil
	}

	eventStream, ok := value.(*pb.EventStream)
	if !ok {
		return &pb.GetEventStreamResponse{Error: "value is not an event stream"}, nil
	}

	return &pb.GetEventStreamResponse{EventStream: eventStream}, nil
}

// GetDataset retrieves dataset details from the specified path
func (s *nexusServer) GetDataset(ctx context.Context, req *pb.GetPathRequest) (*pb.GetDatasetResponse, error) {
	log.Printf("Received dataset request for path: %s\n", req.Path)

	value, err := pathTrie.Get(req.Path)
	if err != nil {
		return &pb.GetDatasetResponse{Error: err.Error()}, nil
	}

	dataset, ok := value.(*pb.Dataset)
	if !ok {
		return &pb.GetDatasetResponse{Error: "value is not a dataset"}, nil
	}

	return &pb.GetDatasetResponse{Dataset: dataset}, nil
}

func main() {
	log.Println("Starting Nexus Server...")

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterNexusServiceServer(s, &nexusServer{})

	log.Printf("Server listening at :%d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
