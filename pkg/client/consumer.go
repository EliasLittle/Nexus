package client

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	pb "nexus/pkg/proto"
)

// GetValue reads a single value from the specified path
/*
func (n *NexusClient) GetValue(path string) (*pb.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	res, err := n.Client.GetValue(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
*/

/*
// GetEventStream subscribes to a Kafka topic and processes events in real-time
func (n *NexusClient) GetEventStream(path string) (<-chan []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	res, err := n.Client.GetEventStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream details: %v", err)
	}

	if res.Error != "" {
		return nil, fmt.Errorf("server error: %s", res.Error)
	}

	eventStream := res.EventStream
	if eventStream == nil {
		return nil, fmt.Errorf("no event stream found at path: %s", path)
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{eventStream.Server}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(eventStream.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to create partition consumer: %v", err)
	}

	messageChan := make(chan []byte)

	go func() {
		defer close(messageChan)
		defer partitionConsumer.Close()
		defer consumer.Close()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				messageChan <- msg.Value
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming message: %v", err)
				return
			}
		}
	}()

	return messageChan, nil
}
*/

// GetDataset reads data from either a file or database table
/*
func (n *NexusClient) GetDataset(path string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetPathRequest{Path: path}
	res, err := n.Client.GetDataset(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset details: %v", err)
	}

	if res.Error != "" {
		return nil, fmt.Errorf("server error: %s", res.Error)
	}

	dataset := res.Dataset
	if dataset == nil {
		return nil, fmt.Errorf("no dataset found at path: %s", path)
	}

	switch d := dataset.Dataset.(type) {
	case *pb.Dataset_IndividualFile:
		return readFile(d.IndividualFile.FilePath)
	case *pb.Dataset_DatabaseTable:
		return queryTable(d.DatabaseTable)
	default:
		return nil, fmt.Errorf("unsupported dataset type")
	}
}
*/

// readFile reads data from a CSV file
func ReadFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// generateConnectionString creates a PostgreSQL connection string from the DatabaseTable
func generateConnectionString(table *pb.DatabaseTable, username, password string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		username, password, table.Host, table.Port, table.DbName)
}

// queryTable reads all data from a database table
func QueryTable(table *pb.DatabaseTable) ([][]string, error) {
	// Get the username from the environment variable
	username := os.Getenv("USER") // For Unix-like systems
	//	if username == "" {
	//		username = os.Getenv("USERNAME") // For Windows systems
	//	}
	if username == "" {
		return nil, fmt.Errorf("could not determine the username from environment variables")
	}

	password := os.Getenv("DB_PASSWORD") // Get the password from the environment variable
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is not set")
	}

	connectionString := generateConnectionString(table, username, password)
	db, err := sql.Open(table.DbType, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s", table.TableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	result := [][]string{columns}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		stringValues := make([]string, len(columns))
		for i, v := range values {
			stringValues[i] = fmt.Sprintf("%v", v)
		}
		result = append(result, stringValues)
	}

	return result, nil
}

// GetData retrieves data from the specified path, returning either a Value or a Dataset
/*
func (n *NexusClient) GetData(path string) (interface{}, error) {
	// First, try to get a Value
	value, err := n.GetValue(path)
	if err == nil {
		return value, nil
	}

	// If getting a Value fails, try to get a Dataset
	dataset, err := n.GetDataset(path)
	if err == nil {
		return dataset, nil
	}

	// If both attempts fail, return the last error
	return nil, fmt.Errorf("failed to get data from path: %s, errors: value error: %v, dataset error: %v", path, err, err)
}
*/
