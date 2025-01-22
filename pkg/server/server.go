package server

import (
	"context"
	"nexus/pkg/logger"
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

// NewServer creates a new NexusServer instance
func NewServer(loadPath string) (*NexusServer, error) {
	log := logger.GetLogger()
	var index *Trie
	var err error

	if loadPath != "" {
		log.Info("Loading index from disk", "path", loadPath)
		index, err = NewTrie(loadPath)
	} else {
		log.Info("No index path provided, creating new Trie")
		index, err = NewTrie()
	}

	if err != nil {
		return nil, err
	}

	return &NexusServer{Index: index}, nil
}

// SaveIndex saves the server's index to disk
func (s *NexusServer) SaveIndex(savePath string) error {
	return s.Index.SaveToDisk(savePath)
}

// RegisterEventStream implements the publisher endpoint for registering event streams
func (s *NexusServer) RegisterEventStream(ctx context.Context, req *pb.RegisterEventStreamRequest) (*pb.RegisterEventStreamResponse, error) {
	log := logger.GetLogger()
	log.Info("Received event stream registration request", "path", req.Path, "topic", req.EventStream.Topic)

	// Insert the event stream into the Trie
	s.Index.Insert(req.Path, req.EventStream) // Store the EventStream
	s.Index.Traverse()                        // Print the Trie after the update

	return &pb.RegisterEventStreamResponse{Success: true}, nil
}

// StoreValue implements the publisher endpoint for storing values
func (s *NexusServer) StoreValue(ctx context.Context, req *pb.StoreValueRequest) (*pb.StoreValueResponse, error) {
	log := logger.GetLogger()
	log.Info("Received value storage request", "path", req.Path)

	// Handle the oneof value types
	switch req.Value.(type) {
	case *pb.StoreValueRequest_StringValue:
		s.Index.Insert(req.Path, req.GetStringValue())
	case *pb.StoreValueRequest_IntValue:
		s.Index.Insert(req.Path, req.GetIntValue())
	case *pb.StoreValueRequest_FloatValue:
		s.Index.Insert(req.Path, req.GetFloatValue())
	default:
		return &pb.StoreValueResponse{Success: false, Error: "invalid value type"}, nil
	}

	//s.Index.Traverse() // Print the Trie after the update
	return &pb.StoreValueResponse{Success: true}, nil
}

// GetNode implements the consumer endpoint for getting node information
func (s *NexusServer) GetNode(ctx context.Context, req *pb.GetPathRequest) (*pb.GetNodeResponse, error) {
	log := logger.GetLogger()
	log.Info("Received request to get node", "path", req.Path)

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
		log.Info("StringValue found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_StringValue{StringValue: v}, ValueType: node.ValueType}, nil
	case *pb.IntValue:
		log.Info("IntValue found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_IntValue{IntValue: v}, ValueType: node.ValueType}, nil
	case *pb.FloatValue:
		log.Info("FloatValue found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_FloatValue{FloatValue: v}, ValueType: node.ValueType}, nil
	case *pb.DatabaseTable:
		log.Info("DatabaseTable found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_DatabaseTable{DatabaseTable: v}, ValueType: node.ValueType}, nil
	case *pb.Directory:
		log.Info("Directory found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_Directory{Directory: v}, ValueType: node.ValueType}, nil
	case *pb.IndividualFile:
		log.Info("IndividualFile found at path", "path", req.Path)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_IndividualFile{IndividualFile: v}, ValueType: node.ValueType}, nil
	case *pb.EventStream:
		log.Info("EventStream found at path", "path", req.Path)
		log.Info("EventStream", "node", v)
		return &pb.GetNodeResponse{Value: &pb.GetNodeResponse_EventStream{EventStream: v}, ValueType: node.ValueType}, nil
	default:
		log.Error("Unknown value type", "type", v)
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
	log := logger.GetLogger()
	log.Info("Received request to list children", "path", req.Path)

	children := s.Index.GetChildren(req.Path)
	return &pb.GetChildrenResponse{Children: children}, nil
}

// RegisterFile implements the publisher endpoint for registering individual files
func (s *NexusServer) RegisterFile(ctx context.Context, req *pb.RegisterFileRequest) (*pb.RegisterFileResponse, error) {
	log := logger.GetLogger()
	log.Info("Received file registration request", "path", req.Path)

	s.Index.Insert(req.Path, req.IndividualFile)
	s.Index.Traverse() // Print the Trie after the update

	return &pb.RegisterFileResponse{Success: true}, nil
}

// RegisterDirectory implements the publisher endpoint for registering directories
func (s *NexusServer) RegisterDirectory(ctx context.Context, req *pb.RegisterDirectoryRequest) (*pb.RegisterDirectoryResponse, error) {
	log := logger.GetLogger()
	log.Info("Received directory registration request", "path", req.Path)

	s.Index.Insert(req.GetPath(), req.GetDirectory())
	s.Index.Traverse() // Print the Trie after the update

	return &pb.RegisterDirectoryResponse{Success: true}, nil
}

// RegisterDatabaseTable implements the publisher endpoint for registering database tables
func (s *NexusServer) RegisterDatabaseTable(ctx context.Context, req *pb.RegisterDatabaseTableRequest) (*pb.RegisterDatabaseTableResponse, error) {
	log := logger.GetLogger()
	log.Info("Received database table registration request", "path", req.Path)

	s.Index.Insert(req.GetPath(), req.GetDatabaseTable())
	s.Index.Traverse() // Print the Trie after the update

	return &pb.RegisterDatabaseTableResponse{Success: true}, nil
}
