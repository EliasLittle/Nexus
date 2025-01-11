package server

import (
	"context"
	"log"
	pb "nexus/pkg/proto"
)

const (
	DefaultPort = 50051
)

// nexusServer implements the NexusService gRPC service
type NexusServer struct {
	pb.UnimplementedNexusServiceServer
	Index *Trie
}

// RegisterEventStream implements the publisher endpoint for registering event streams
func (s *NexusServer) RegisterEventStream(ctx context.Context, req *pb.RegisterEventStreamRequest) (*pb.RegisterEventStreamResponse, error) {
	log.Printf("Received event stream registration request for path: %s, topic: %s\n", req.Path, req.EventStream.Topic)

	// Insert the event stream into the Trie
	s.Index.Insert(req.Path, req.EventStream) // Store the EventStream
	s.Index.Traverse()                        // Print the Trie after the update

	return &pb.RegisterEventStreamResponse{Success: true}, nil
}

// RegisterDataset implements the publisher endpoint for registering datasets
func (s *NexusServer) RegisterDataset(ctx context.Context, req *pb.RegisterDatasetRequest) (*pb.RegisterDatasetResponse, error) {
	log.Printf("Received dataset registration request for path: %s\n", req.Path)

	// Insert the dataset into the Trie
	s.Index.Insert(req.Path, req.Dataset) // Store the Dataset
	s.Index.Traverse()                    // Print the Trie after the update

	return &pb.RegisterDatasetResponse{Success: true}, nil
}

// StoreValue implements the publisher endpoint for storing values
func (s *NexusServer) StoreValue(ctx context.Context, req *pb.StoreValueRequest) (*pb.StoreValueResponse, error) {
	log.Printf("Received value storage request for path: %s\n", req.Path)

	s.Index.Insert(req.Path, req.DirectValue)

	//s.Index.Traverse() // Print the Trie after the update
	return &pb.StoreValueResponse{Success: true}, nil
}

// GetValue implements the consumer endpoint for retrieving values
func (s *NexusServer) GetValue(ctx context.Context, req *pb.GetPathRequest) (*pb.GetValueResponse, error) {
	log.Printf("Received value request for path: %s\n", req.Path)

	value, err := s.Index.Get(req.Path)

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

// GetEventStream retrieves event stream details from the specified path
func (s *NexusServer) GetEventStream(ctx context.Context, req *pb.GetPathRequest) (*pb.GetEventStreamResponse, error) {
	log.Printf("Received event stream request for path: %s\n", req.Path)

	value, err := s.Index.Get(req.Path)
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
func (s *NexusServer) GetDataset(ctx context.Context, req *pb.GetPathRequest) (*pb.GetDatasetResponse, error) {
	log.Printf("Received dataset request for path: %s\n", req.Path)

	value, err := s.Index.Get(req.Path)
	if err != nil {
		return &pb.GetDatasetResponse{Error: err.Error()}, nil
	}

	dataset, ok := value.(*pb.Dataset)
	if !ok {
		return &pb.GetDatasetResponse{Error: "value is not a dataset"}, nil
	}

	return &pb.GetDatasetResponse{Dataset: dataset}, nil
}

func (s *NexusServer) GetPathType(ctx context.Context, req *pb.GetPathRequest) (*pb.GetPathTypeResponse, error) {
	pathType, err := s.Index.GetType(req.Path)
	if err != nil {
		return &pb.GetPathTypeResponse{Error: err.Error()}, nil
	}

	return &pb.GetPathTypeResponse{PathType: pathType}, nil
}

// ListChildren implements the endpoint for listing children nodes
func (s *NexusServer) GetChildren(ctx context.Context, req *pb.GetChildrenRequest) (*pb.GetChildrenResponse, error) {
	log.Printf("Received request to list children for path: %s\n", req.Path)

	children := s.Index.GetChildren(req.Path)

	return &pb.GetChildrenResponse{Children: children}, nil
}
