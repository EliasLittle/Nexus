package client

import (
	"context"
	"time"

	pb "nexus/pkg/proto"

	"google.golang.org/grpc"
)

// PublishEventStream registers an event stream with the Nexus server
func PublishEventStream(conn *grpc.ClientConn, path string, eventStream *pb.EventStream) error {
	client := pb.NewNexusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterEventStreamRequest{
		Path:        path,
		EventStream: eventStream,
	}

	_, err := client.RegisterEventStream(ctx, req)
	return err
}

// PublishDataset registers a dataset with the Nexus server
func PublishDataset(conn *grpc.ClientConn, path string, dataset *pb.Dataset) error {
	client := pb.NewNexusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterDatasetRequest{
		Path:    path,
		Dataset: dataset,
	}

	_, err := client.RegisterDataset(ctx, req)
	return err
}

// PublishValue stores a value directly to the Nexus server
func PublishValue(conn *grpc.ClientConn, path string, directValue *pb.DirectValue) error {
	client := pb.NewNexusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.StoreValueRequest{
		Path:        path,
		DirectValue: directValue,
	}

	_, err := client.StoreValue(ctx, req)
	return err
}

/*
func publishMain() {
	// Connect to the Nexus server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Example usage of publishing an event stream
	eventStream := &pb.EventStream{
		Server: "localhost:9092",
		Topic:  "sensor_data",
	}
	err = PublishEventStream(conn, "/events/sensors", eventStream)
	if err != nil {
		log.Fatalf("Failed to publish event stream: %v", err)
	}

	// Example usage of publishing a dataset
	dataset := &pb.Dataset{
		DatasetType: &pb.Dataset_IndividualFile{
			IndividualFile: &pb.IndividualFile{
				FileType:    "csv",
				FilePath:    "path/to/users.csv",
				ColumnNames: []string{"id", "name", "email"},
			},
		},
	}
	err = PublishDataset(conn, "/datasets/users", dataset)
	if err != nil {
		log.Fatalf("Failed to publish dataset: %v", err)
	}

	// Example usage of publishing a value
	directValue := &pb.DirectValue{
		DataStructure: "string",
		Value:         "operational",
	}
	err = PublishValue(conn, "/values/system-status", directValue)
	if err != nil {
		log.Fatalf("Failed to publish value: %v", err)
	}

	log.Println("Data published successfully!")
}
*/
