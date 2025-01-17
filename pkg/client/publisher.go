package client

import (
	"context"
	"time"

	pb "nexus/pkg/proto"
)

// PublishEventStream registers an event stream with the Nexus server
func (n *NexusClient) PublishEventStream(path string, eventStream *pb.EventStream) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterEventStreamRequest{
		Path:        path,
		EventStream: eventStream,
	}

	_, err := n.Client.RegisterEventStream(ctx, req)
	return err
}

// PublishDataset registers a dataset with the Nexus server
func (n *NexusClient) PublishIndividualFile(path string, file *pb.IndividualFile) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterFileRequest{
		Path:           path,
		IndividualFile: file,
	}

	_, err := n.Client.RegisterFile(ctx, req)
	return err
}

func (n *NexusClient) PublishDirectory(path string, directory *pb.Directory) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDirectoryRequest{
		Path:      path,
		Directory: directory,
	}
	_, err := n.Client.RegisterDirectory(ctx, req)
	return err
}

func (n *NexusClient) PublishDatabaseTable(path string, table *pb.DatabaseTable) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDatabaseTableRequest{
		Path:          path,
		DatabaseTable: table,
	}
	_, err := n.Client.RegisterDatabaseTable(ctx, req)
	return err
}

// PublishValue stores a value directly to the Nexus server
func (n *NexusClient) PublishValue(path string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req *pb.StoreValueRequest
	switch v := value.(type) {
	case string:
		req = &pb.StoreValueRequest{
			Path:  path,
			Value: &pb.StoreValueRequest_StringValue{StringValue: &pb.StringValue{Value: v}},
		}
	case int:
		req = &pb.StoreValueRequest{
			Path:  path,
			Value: &pb.StoreValueRequest_IntValue{IntValue: &pb.IntValue{Value: int32(v)}},
		}
	case float64:
		req = &pb.StoreValueRequest{
			Path:  path,
			Value: &pb.StoreValueRequest_FloatValue{FloatValue: &pb.FloatValue{Value: float32(v)}},
		}
	default:
		return nil // or handle the error as needed
	}

	_, err := n.Client.StoreValue(ctx, req)
	return err
}
