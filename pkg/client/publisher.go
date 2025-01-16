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
func (n *NexusClient) PublishDataset(path string, dataset *pb.Dataset) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDatasetRequest{
		Path:    path,
		Dataset: dataset,
	}

	_, err := n.Client.RegisterDataset(ctx, req)
	return err
}

// PublishValue stores a value directly to the Nexus server
func (n *NexusClient) PublishValue(path string, value *pb.Value) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.StoreValueRequest{
		Path:  path,
		Value: value,
	}

	_, err := n.Client.StoreValue(ctx, req)
	return err
}
