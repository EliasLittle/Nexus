package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "nexus/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultConnection = "localhost:50051" // Default connection string

func createGRPCConnection(connStr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createDirectValue(value string) *pb.DirectValue {
	return &pb.DirectValue{
		Value: &pb.DirectValue_StringValue{StringValue: value},
	}
}

func createEventStream(topic string) *pb.EventStream {
	return &pb.EventStream{
		Server: "localhost:9092",
		Topic:  topic,
	}
}

func createDataset(filePath string) *pb.Dataset {
	return &pb.Dataset{
		Dataset: &pb.Dataset_IndividualFile{
			IndividualFile: &pb.IndividualFile{
				FileType:    "csv",
				FilePath:    filePath,
				ColumnNames: []string{"id", "name", "email"},
			},
		},
	}
}

func listChildren(conn *grpc.ClientConn, path string) {
	client := pb.NewNexusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ListChildrenRequest{Path: path}
	resp, err := client.ListChildren(ctx, req)
	if err != nil {
		fmt.Println("Error listing children:", err)
		return
	}

	fmt.Printf("Children of path '%s':\n", path)
	for _, child := range resp.Children {
		fmt.Println(child)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nexus-client publish|consume|list")
		os.Exit(1)
	}

	conn, err := createGRPCConnection(defaultConnection)
	if err != nil {
		fmt.Println("Failed to create gRPC connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if len(os.Args) > 4 && strings.HasPrefix(os.Args[4], "conn=") {
		conn, err = createGRPCConnection(strings.TrimPrefix(os.Args[4], "conn="))
		if err != nil {
			fmt.Println("Failed to create gRPC connection with conn=", strings.TrimPrefix(os.Args[4], "conn="), " :", err)
			os.Exit(1)
		}
	}

	switch os.Args[1] {
	case "publish":
		if len(os.Args) < 4 {
			fmt.Println("Expected 'dataset', 'event', or 'value' as argument for publish, followed by a path")
			os.Exit(1)
		}
		path := os.Args[3]
		switch os.Args[2] {
		case "dataset":
			dataset := createDataset(os.Args[4])
			err = PublishDataset(conn, path, dataset)
			if err != nil {
				fmt.Println("Failed to publish dataset:", err)
				os.Exit(1)
			}
		case "event":
			eventStream := createEventStream(os.Args[4])
			err = PublishEventStream(conn, path, eventStream)
			if err != nil {
				fmt.Println("Failed to publish event stream:", err)
				os.Exit(1)
			}
		case "value":
			directValue := createDirectValue(os.Args[4])
			err = PublishValue(conn, path, directValue)
			if err != nil {
				fmt.Println("Failed to publish value:", err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unknown publish type. Use 'dataset', 'event', or 'value'.")
			os.Exit(1)
		}
	case "consume":
		var path string
		if len(os.Args) < 3 {
			path = "/" // Default to root if no path is provided
		} else {
			path = os.Args[3]
		}
		switch os.Args[2] {
		case "value":
			value, err := ConsumeValue(conn, path)
			if err != nil {
				fmt.Println("Failed to consume value:", err)
				os.Exit(1)
			}
			fmt.Printf("Consumed value: %v\n", value)
		case "dataset":
			fmt.Println("Consuming dataset...")
			dataset, err := ConsumeDataset(conn, path)
			if err != nil {
				fmt.Println("Failed to consume dataset:", err)
				os.Exit(1)
			}
			fmt.Printf("Consumed dataset: %v\n", dataset)
		case "event":
			eventStream, err := ConsumeEventStream(conn, path)
			if err != nil {
				fmt.Println("Failed to consume event stream:", err)
				os.Exit(1)
			}
			fmt.Printf("Consumed event stream: %v\n", eventStream)
		default:
			fmt.Println("Unknown consume type. Use 'value', 'dataset', or 'event'.")
			os.Exit(1)
		}
	case "list":
		var path string
		if len(os.Args) > 2 {
			path = os.Args[2]
		} else {
			path = "/" // Default to root if no path is provided
		}
		listChildren(conn, path)
	default:
		fmt.Println("Unknown command. Use 'publish' or 'consume'.")
		os.Exit(1)
	}
}
