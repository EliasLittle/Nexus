package client

import (
	"context"
	"fmt"
	"time"

	pb "nexus/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DefaultConnection = "localhost:50051" // Default connection string

type NexusClient struct {
	Client pb.NexusServiceClient
}

func NewNexusClient(conn *grpc.ClientConn) *NexusClient {
	return &NexusClient{
		Client: pb.NewNexusServiceClient(conn),
	}
}

func CreateGRPCConnection(connStr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateDirectValue(value string) *pb.DirectValue {
	return &pb.DirectValue{
		Value: &pb.DirectValue_StringValue{StringValue: value},
	}
}

func CreateEventStream(topic string) *pb.EventStream {
	return &pb.EventStream{
		Server: "localhost:9092",
		Topic:  topic,
	}
}

func CreateDataset(filePath string) *pb.Dataset {
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

func (n *NexusClient) GetChildren(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetChildrenRequest{Path: path}
	resp, err := n.Client.GetChildren(ctx, req)
	if err != nil {
		fmt.Println("Error listing children:", err)
		return nil, err
	}

	return resp.Children, nil
}

func (n *NexusClient) GetPathType(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	resp, err := n.Client.GetPathType(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.PathType, nil
}
