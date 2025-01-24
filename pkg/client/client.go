package client

import (
	"context"
	"fmt"
	"nexus/pkg/logger"
	"os"
	"path"
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
	log := logger.GetLogger()
	log.Debug("Creating gRPC connection", "address", connStr)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to create gRPC connection", "error", err)
		return nil, err
	}
	return conn, nil
}

func CreateEventStream(topic string) *pb.EventStream {
	log := logger.GetLogger()
	log.Debug("Creating event stream", "topic", topic)
	return &pb.EventStream{
		Server: "localhost:9092",
		Topic:  topic,
	}
}

func CreateIndividualFile(filePath string) *pb.IndividualFile {
	log := logger.GetLogger()
	fileType := path.Ext(filePath) // Get the file extension
	log.Debug("Creating individual file", "path", filePath, "type", fileType)
	return &pb.IndividualFile{
		FileType:    fileType,
		FilePath:    filePath,
		ColumnNames: []string{},
	}
}

func CreateDirectory(directoryPath string) (*pb.Directory, error) {
	log := logger.GetLogger()
	log.Debug("Creating directory", "path", directoryPath)

	// Get the file type from the first file in the directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Error("Failed to read directory", "error", err)
		return nil, err
	}

	var fileType string
	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
			if fileType == "" {
				fileType = path.Ext(file.Name())
			}
		}
	}

	if fileType == "" {
		log.Error("No files found in directory")
		return nil, fmt.Errorf("no files found in directory")
	}

	log.Debug("Directory created", "type", fileType, "count", fileCount)
	return &pb.Directory{
		FileType:      fileType,
		DirectoryPath: directoryPath,
		FileCount:     int32(fileCount),
	}, nil
}

func CreateDatabaseTable(dbType string, host string, port int32, dbName string, tableName string) *pb.DatabaseTable {
	log := logger.GetLogger()
	log.Debug("Creating database table", "type", dbType, "host", host, "port", port, "db", dbName, "table", tableName)
	return &pb.DatabaseTable{
		DbType:    dbType,
		Host:      host,
		Port:      port,
		DbName:    dbName,
		TableName: tableName,
	}
}

func (n *NexusClient) GetChildren(path string) ([]*pb.ChildInfo, error) {
	log := logger.GetLogger()
	log.Debug("Getting children", "path", path)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetChildrenRequest{Path: path}
	res, err := n.Client.GetChildren(ctx, req)
	if err != nil {
		log.Error("Failed to get children", "error", err)
		return nil, err
	}

	log.Debug("Got children", "path", path, "children", res.Children)
	return res.Children, nil
}

// TODO: Add rpc to get path type to reduce unneeded data transfer
func (n *NexusClient) GetPathType(path string) (string, error) {
	log := logger.GetLogger()
	log.Debug("Getting path type", "path", path)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	res, err := n.Client.GetNode(ctx, req)
	if err != nil {
		log.Error("Failed to get path type", "error", err)
		return "", err
	}

	log.Debug("Got path type", "path", path, "type", res.ValueType)
	return res.ValueType, nil
}

// TODO: Add rpc to get path data to reduce unneeded data transfer
func (n *NexusClient) Get(path string) (interface{}, error) {
	log := logger.GetLogger()
	log.Debug("Getting value", "path", path)

	value, _, err := n.GetFull(path)
	if err != nil {
		log.Error("Failed to get value", "error", err)
		return nil, err
	}

	log.Debug("Got value", "path", path, "value", value)
	return value, nil
}

func (n *NexusClient) GetFull(path string) (interface{}, string, error) {
	log := logger.GetLogger()
	log.Debug("Getting full value", "path", path)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	node, err := n.Client.GetNode(ctx, req)
	if err != nil {
		log.Error("Failed to get full value", "error", err)
		return nil, "", err
	}

	switch valType := node.GetValueType(); valType {
	case "StringValue":
		return node.GetStringValue(), valType, nil
	case "IntValue":
		return node.GetIntValue(), valType, nil
	case "FloatValue":
		return node.GetFloatValue(), valType, nil
	case "DatabaseTable":
		return node.GetDatabaseTable(), valType, nil
	case "Directory":
		return node.GetDirectory(), valType, nil
	case "IndividualFile":
		return node.GetIndividualFile(), valType, nil
	case "EventStream":
		log.Info("EventStream found", "node", node)
		return node.GetEventStream(), valType, nil
	default:
		log.Error("Unknown value type", "type", valType)
		return nil, "", fmt.Errorf("unknown value type: %s", valType)
	}
}
