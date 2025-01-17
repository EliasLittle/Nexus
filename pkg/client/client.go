package client

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

func CreateEventStream(topic string) *pb.EventStream {
	return &pb.EventStream{
		Server: "localhost:9092",
		Topic:  topic,
	}
}

func CreateIndividualFile(filePath string) *pb.IndividualFile {
	fileType := path.Ext(filePath) // Get the file extension
	return &pb.IndividualFile{
		FileType:    fileType,
		FilePath:    filePath,
		ColumnNames: []string{},
	}
}

func CreateDirectory(directoryPath string) *pb.Directory {
	files, err := os.ReadDir(directoryPath) // Read the directory contents
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	var fileType string
	fileCount := 0
	extensionSet := make(map[string]struct{}) // To track unique file extensions

	for _, file := range files {
		if !file.IsDir() { // Only consider files, not directories
			fileCount++
			ext := filepath.Ext(file.Name()) // Get the file extension
			extensionSet[ext] = struct{}{}   // Add the extension to the set
			if fileType == "" {
				fileType = ext // Set fileType to the first file's extension
			}
		}
	}

	// Check if all files have the same extension
	if len(extensionSet) > 1 {
		fmt.Println("Not all files have the same extension.")
		return nil
	}

	return &pb.Directory{
		FileType:      fileType,
		DirectoryPath: directoryPath,
		FileCount:     int32(fileCount),
	}
}

func CreateDatabaseTable(dbType string, host string, port int32, dbName string, tableName string) *pb.DatabaseTable {
	return &pb.DatabaseTable{
		DbType:    dbType,
		Host:      host,
		Port:      port,
		DbName:    dbName,
		TableName: tableName,
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

// TODO: Add rpc to get path type to reduce unneeded data transfer
func (n *NexusClient) GetPathType(path string) (string, error) {
	_, valueType, err := n.GetFull(path)
	if err != nil {
		return "", err
	}

	return valueType, nil
}

// TODO: Add rpc to get path data to reduce unneeded data transfer
func (n *NexusClient) Get(path string) (interface{}, error) {
	value, _, err := n.GetFull(path)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("client.go Get| Value %v has type: %T", value, value)
	return value, nil
}

func (n *NexusClient) GetFull(path string) (interface{}, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	node, err := n.Client.GetNode(ctx, req)
	if err != nil {
		return nil, "", err
	}

	switch node.ValueType {
	case "StringValue":
		return node.GetStringValue(), node.ValueType, nil
	case "IntValue":
		return node.GetIntValue(), node.ValueType, nil
	case "FloatValue":
		return node.GetFloatValue(), node.ValueType, nil
	case "DatabaseTable":
		return node.GetDatabaseTable(), node.ValueType, nil
	case "Directory":
		return node.GetDirectory(), node.ValueType, nil
	case "IndividualFile":
		return node.GetIndividualFile(), node.ValueType, nil
	case "EventStream":
		return node.GetEventStream(), node.ValueType, nil
	default:
		return nil, "", fmt.Errorf("unsupported value type: %v at path: %s", node.ValueType, path)
	}
}
