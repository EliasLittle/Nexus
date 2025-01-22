package client

import (
	"context"
	"fmt"
	"nexus/pkg/logger"
	pb "nexus/pkg/proto"
	"time"
)

// PublishEventStream registers an event stream with the Nexus server
func (n *NexusClient) PublishEventStream(path string, eventStream *pb.EventStream) error {
	log := logger.GetLogger()
	log.Debug("Publishing event stream", "path", path, "topic", eventStream.Topic)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterEventStreamRequest{
		Path:        path,
		EventStream: eventStream,
	}

	_, err := n.Client.RegisterEventStream(ctx, req)
	if err != nil {
		log.Error("Failed to publish event stream", "error", err)
		return err
	}

	log.Debug("Event stream published successfully", "path", path)
	return nil
}

// PublishIndividualFile registers a file with the Nexus server
func (n *NexusClient) PublishIndividualFile(path string, file *pb.IndividualFile) error {
	log := logger.GetLogger()
	log.Debug("Publishing individual file", "path", path, "file", file.FilePath)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterFileRequest{
		Path:           path,
		IndividualFile: file,
	}

	_, err := n.Client.RegisterFile(ctx, req)
	if err != nil {
		log.Error("Failed to publish individual file", "error", err)
		return err
	}

	log.Debug("Individual file published successfully", "path", path)
	return nil
}

// PublishDirectory registers a directory with the Nexus server
func (n *NexusClient) PublishDirectory(path string, directory *pb.Directory) error {
	log := logger.GetLogger()
	log.Debug("Publishing directory", "path", path, "directory", directory.DirectoryPath)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDirectoryRequest{
		Path:      path,
		Directory: directory,
	}

	_, err := n.Client.RegisterDirectory(ctx, req)
	if err != nil {
		log.Error("Failed to publish directory", "error", err)
		return err
	}

	log.Debug("Directory published successfully", "path", path)
	return nil
}

// PublishDatabaseTable registers a database table with the Nexus server
func (n *NexusClient) PublishDatabaseTable(path string, table *pb.DatabaseTable) error {
	log := logger.GetLogger()
	log.Debug("Publishing database table", "path", path, "table", table.TableName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDatabaseTableRequest{
		Path:          path,
		DatabaseTable: table,
	}

	_, err := n.Client.RegisterDatabaseTable(ctx, req)
	if err != nil {
		log.Error("Failed to publish database table", "error", err)
		return err
	}

	log.Debug("Database table published successfully", "path", path)
	return nil
}

// PublishValue publishes a value to the Nexus server
func (n *NexusClient) PublishValue(path string, value interface{}) error {
	log := logger.GetLogger()
	log.Debug("Publishing value", "path", path, "value", value)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req *pb.StoreValueRequest
	switch v := value.(type) {
	case string:
		req = &pb.StoreValueRequest{
			Path: path,
			Value: &pb.StoreValueRequest_StringValue{
				StringValue: &pb.StringValue{Value: v},
			},
		}
	case int32:
		req = &pb.StoreValueRequest{
			Path: path,
			Value: &pb.StoreValueRequest_IntValue{
				IntValue: &pb.IntValue{Value: v},
			},
		}
	case float32:
		req = &pb.StoreValueRequest{
			Path: path,
			Value: &pb.StoreValueRequest_FloatValue{
				FloatValue: &pb.FloatValue{Value: v},
			},
		}
	default:
		err := fmt.Errorf("unsupported value type: %T", value)
		log.Error("Failed to publish value", "error", err)
		return err
	}

	_, err := n.Client.StoreValue(ctx, req)
	if err != nil {
		log.Error("Failed to publish value", "error", err)
		return err
	}

	log.Debug("Value published successfully", "path", path)
	return nil
}
