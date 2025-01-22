package client

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"nexus/pkg/logger"
	"os"

	pb "nexus/pkg/proto"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
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

// GetEventStream subscribes to a Kafka topic and processes events in real-time
func GetEventStream(es *pb.EventStream) (<-chan []byte, error) {
	log := logger.GetLogger()
	if es == nil {
		log.Error("No event stream provided")
		return nil, fmt.Errorf("no event stream provided")
	}

	log.Debug("Creating Kafka consumer", "server", es.Server, "topic", es.Topic)
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{es.Server}, config)
	if err != nil {
		log.Error("Failed to create Kafka consumer", "error", err)
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(es.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.Close()
		log.Error("Failed to create partition consumer", "error", err)
		return nil, fmt.Errorf("failed to create partition consumer: %v", err)
	}

	// Create a channel to receive messages
	messages := make(chan []byte)

	// Start a goroutine to handle messages
	go func() {
		defer func() {
			partitionConsumer.Close()
			consumer.Close()
			close(messages)
		}()

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				messages <- msg.Value
			case err := <-partitionConsumer.Errors():
				log.Error("Error consuming message", "error", err)
			}
		}
	}()

	log.Debug("Event stream consumer started", "server", es.Server, "topic", es.Topic)
	return messages, nil
}

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

// ReadFile reads data from a file and returns it as a slice of string slices
func ReadFile(filePath string) ([][]string, error) {
	log := logger.GetLogger()
	log.Debug("Reading file", "path", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Error("Failed to open file", "error", err)
		return nil, err
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
			log.Error("Failed to read record", "error", err)
			return nil, err
		}
		records = append(records, record)
	}

	log.Debug("File read successfully", "path", filePath, "records", len(records))
	return records, nil
}

// generateConnectionString creates a PostgreSQL connection string from the DatabaseTable
func generateConnectionString(table *pb.DatabaseTable, username, password string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		table.Host, table.Port, username, password, table.DbName)
}

// QueryTable executes a query on a database table and returns the results
func QueryTable(table *pb.DatabaseTable) ([][]string, error) {
	log := logger.GetLogger()
	log.Debug("Querying table", "type", table.DbType, "host", table.Host, "db", table.DbName, "table", table.TableName)
	// Get the username from the environment variable
	username := os.Getenv("USER") // For Unix-like systems
	fmt.Println("using username: ", username)
	//	if username == "" {
	//		username = os.Getenv("USERNAME") // For Windows systems
	//	}
	if username == "" {
		return nil, fmt.Errorf("could not determine the username from environment variables")
	}

	password := os.Getenv("DB_PASSWORD") // Get the password from the environment variable
	fmt.Println("using password: ", password)
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is not set")
	}

	connStr := generateConnectionString(table, username, password)
	db, err := sql.Open(table.DbType, connStr)
	if err != nil {
		log.Error("Failed to open database connection", "error", err)
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM %s limit 100", table.TableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Error("Failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Error("Failed to get column names", "error", err)
		return nil, err
	}

	// Prepare result slice with column names as first row
	result := [][]string{columns}

	// Prepare value holders
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Iterate through rows
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Error("Failed to scan row", "error", err)
			return nil, err
		}

		// Convert values to strings
		var record []string
		for _, val := range values {
			record = append(record, fmt.Sprintf("%v", val))
		}
		result = append(result, record)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error iterating rows", "error", err)
		return nil, err
	}

	log.Debug("Query completed successfully", "rows", len(result)-1)
	return result, nil
}
