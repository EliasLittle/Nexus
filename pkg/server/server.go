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
/*
func (s *NexusServer) RegisterDataset(ctx context.Context, req *pb.RegisterDatasetRequest) (*pb.RegisterDatasetResponse, error) {
	log.Printf("Received dataset registration request for path: %s\n", req.Path)

	// Insert the dataset into the Trie
	s.Index.Insert(req.Path, req.Dataset) // Store the Dataset
	s.Index.Traverse()                    // Print the Trie after the update

	return &pb.RegisterDatasetResponse{Success: true}, nil
}
*/

// StoreValue implements the publisher endpoint for storing values
func (s *NexusServer) StoreValue(ctx context.Context, req *pb.StoreValueRequest) (*pb.StoreValueResponse, error) {
	log.Printf("Received value storage request for path: %s\n", req.Path)
	//log.Printf("req.Value has type: %s", reflect.TypeOf(req.Value).Name())

	// Handle the oneof value types
	switch req.Value.(type) {
	case *pb.StoreValueRequest_StringValue:
		//log.Printf("Value has type: %s", reflect.TypeOf(*req.GetStringValue()).Name())
		s.Index.Insert(req.Path, req.GetStringValue())
	case *pb.StoreValueRequest_IntValue:
		//log.Printf("Value has type: %s", reflect.TypeOf(*req.GetIntValue()).Name())
		s.Index.Insert(req.Path, req.GetIntValue())
	case *pb.StoreValueRequest_FloatValue:
		//log.Printf("Value has type: %s", reflect.TypeOf(*req.GetFloatValue()).Name())
		s.Index.Insert(req.Path, req.GetFloatValue())
	default:
		return &pb.StoreValueResponse{Success: false, Error: "invalid value type"}, nil
	}

	//s.Index.Traverse() // Print the Trie after the update
	return &pb.StoreValueResponse{Success: true}, nil
}

func (s *NexusServer) GetNode(ctx context.Context, req *pb.GetPathRequest) (*pb.GetNodeResponse, error) {
	log.Printf("Received request to get node for path: %s\n", req.Path)
	node, err := s.Index.GetNode(req.Path)
	if err != nil {
		return &pb.GetNodeResponse{Error: err.Error()}, nil
	}

	if node.ValueType == "InteriorNode" {
		log.Printf("InteriorNode at path: %s\n", req.Path)
		return &pb.GetNodeResponse{Value: nil, ValueType: node.ValueType}, nil
	}

	switch v := node.Value.(type) {
	case *pb.StringValue:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_StringValue{StringValue: v}, ValueType: node.ValueType}, nil
	case *pb.IntValue:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_IntValue{IntValue: v}, ValueType: node.ValueType}, nil
	case *pb.FloatValue:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_FloatValue{FloatValue: v}, ValueType: node.ValueType}, nil
	case *pb.DatabaseTable:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_DatabaseTable{DatabaseTable: v}, ValueType: node.ValueType}, nil
	case *pb.Directory:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_Directory{Directory: v}, ValueType: node.ValueType}, nil
	case *pb.IndividualFile:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_IndividualFile{IndividualFile: v}, ValueType: node.ValueType}, nil
	case *pb.EventStream:
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_EventStream{EventStream: v}, ValueType: node.ValueType}, nil
	default:
		log.Printf("Unknown value type: %v\n", v)
		return &pb.GetNodeResponse{Error: "unknown value type"}, nil
	}
}

// GetValue implements the consumer endpoint for retrieving values
/*
func (s *NexusServer) GetValue(ctx context.Context, req *pb.GetPathRequest) (*pb.Value, error) {
	log.Printf("Received value request for path: %s\n", req.Path)

	node, err := s.Index.GetNode(req.Path)
	log.Printf("Node: %v", node)

	if err != nil {
		return &pb.Value{Error: err.Error()}, nil
	}

	switch v := node.Value.(type) {
	case *pb.StringValue:
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: v}}, nil
	case *pb.IntValue:
		return &pb.Value{Value: &pb.Value_IntValue{IntValue: v}}, nil
	case *pb.FloatValue:
		return &pb.Value{Value: &pb.Value_FloatValue{FloatValue: v}}, nil
	default:
		log.Printf("Value is not of type StringValue, IntValue, or FloatValue: %v", v)
		return &pb.Value{Error: "value is not of type StringValue, IntValue, or FloatValue"}, nil
	}
}
*/

// GetEventStream retrieves event stream details from the specified path
/*
func (s *NexusServer) GetEventStream(ctx context.Context, req *pb.GetPathRequest) (*pb.GetEventStreamResponse, error) {
	log.Printf("Received event stream request for path: %s\n", req.Path)

	node, err := s.Index.GetNode(req.Path)
	if err != nil {
		return &pb.GetEventStreamResponse{Error: err.Error()}, nil
	}

	eventStream, ok := node.Value.(*pb.EventStream)
	if !ok {
		return &pb.GetEventStreamResponse{Error: "value is not an event stream"}, nil
	}

	return &pb.GetEventStreamResponse{EventStream: eventStream}, nil
}
*/

// GetDataset retrieves dataset details from the specified path
/*
func (s *NexusServer) GetDataset(ctx context.Context, req *pb.GetPathRequest) (*pb.GetDatasetResponse, error) {
	log.Printf("Received dataset request for path: %s\n", req.Path)

	node, err := s.Index.GetNode(req.Path)
	if err != nil {
		return &pb.GetDatasetResponse{Error: err.Error()}, nil
	}

	dataset, ok := node.Value.(*pb.Dataset)
	if !ok {
		return &pb.GetDatasetResponse{Error: "value is not a dataset"}, nil
	}

	return &pb.GetDatasetResponse{Dataset: dataset}, nil
}
*/

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

// RegisterFile implements the publisher endpoint for registering files
func (s *NexusServer) RegisterFile(ctx context.Context, req *pb.RegisterFileRequest) (*pb.RegisterFileResponse, error) {
	log.Printf("Received file registration request for path: %s\n", req.Path)

	// Insert the file into the Trie
	//log.Printf("File has type: %s", reflect.TypeOf(*req.GetIndividualFile()).Name())
	s.Index.Insert(req.Path, req.GetIndividualFile()) // Store the File
	s.Index.Traverse()                                // Print the Trie after the update

	return &pb.RegisterFileResponse{Success: true}, nil
}

// RegisterDirectory implements the publisher endpoint for registering directories
func (s *NexusServer) RegisterDirectory(ctx context.Context, req *pb.RegisterDirectoryRequest) (*pb.RegisterDirectoryResponse, error) {
	log.Printf("Received directory registration request for path: %s\n", req.Path)

	// Insert the directory into the Trie
	s.Index.Insert(req.Path, req.GetDirectory()) // Store the Directory
	s.Index.Traverse()                           // Print the Trie after the update

	return &pb.RegisterDirectoryResponse{Success: true}, nil
}

// RegisterDatabaseTable implements the publisher endpoint for registering database tables
func (s *NexusServer) RegisterDatabaseTable(ctx context.Context, req *pb.RegisterDatabaseTableRequest) (*pb.RegisterDatabaseTableResponse, error) {
	log.Printf("Received database table registration request for path: %s\n", req.Path)

	// Insert the database table into the Trie
	s.Index.Insert(req.Path, req.GetDatabaseTable()) // Store the DatabaseTable
	s.Index.Traverse()                               // Print the Trie after the update

	return &pb.RegisterDatabaseTableResponse{Success: true}, nil
}
