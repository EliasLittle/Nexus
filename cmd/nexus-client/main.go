package main

import (
	"fmt"
	"os"
	"strings"

	nc "nexus/pkg/client"
	pb "nexus/pkg/proto"
	"strconv"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: nexus-client publish|consume|list")
		os.Exit(1)
	}

	conn, err := nc.CreateGRPCConnection(nc.DefaultConnection)
	if err != nil {
		fmt.Println("Failed to create gRPC connection:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if len(os.Args) > 4 && strings.HasPrefix(os.Args[4], "conn=") {
		conn, err = nc.CreateGRPCConnection(strings.TrimPrefix(os.Args[4], "conn="))
		if err != nil {
			fmt.Println("Failed to create gRPC connection with conn=", strings.TrimPrefix(os.Args[4], "conn="), " :", err)
			os.Exit(1)
		}
	}

	client := nc.NewNexusClient(conn)

	switch os.Args[1] {
	case "publish":
		if len(os.Args) < 4 {
			fmt.Println(`
			Usage: nexus-client publish <type> <path> <data>

			Types:
			file - Publish a file
			directory - Publish a directory
			DBTable - Publish a database table
			event - Publish an event stream
			value - Publish a value
			`)
			os.Exit(1)
		}
		path := os.Args[3]
		switch os.Args[2] {
		case "file":
			file := nc.CreateIndividualFile(os.Args[4])
			err = client.PublishIndividualFile(path, file)
			if err != nil {
				fmt.Println("Failed to publish dataset:", err)
				os.Exit(1)
			}
		case "directory":
			directory, err := nc.CreateDirectory(os.Args[4])
			if err != nil {
				fmt.Println("Failed to create directory:", err)
				os.Exit(1)
			}
			err = client.PublishDirectory(path, directory)
			if err != nil {
				fmt.Println("Failed to publish directory:", err)
				os.Exit(1)
			}
		case "DBTable":
			if len(os.Args) < 8 {
				fmt.Println(`Usage: nexus-client publish DBTable <path> <db_type> <host> <port> <db_name> <table_name>`)
				os.Exit(1)
			}
			port, err := strconv.Atoi(os.Args[6])
			if err != nil {
				fmt.Println("Failed to convert port to int:", err)
				os.Exit(1)
			}
			database := nc.CreateDatabaseTable(os.Args[4], os.Args[5], int32(port), os.Args[7], os.Args[8])
			err = client.PublishDatabaseTable(path, database)
			if err != nil {
				fmt.Println("Failed to publish database table:", err)
				os.Exit(1)
			}
		case "event":
			eventStream := nc.CreateEventStream(os.Args[4])
			if len(os.Args) > 5 {
				eventStream.Server = os.Args[5]
			}
			err = client.PublishEventStream(path, eventStream)
			if err != nil {
				fmt.Println("Failed to publish event stream:", err)
				os.Exit(1)
			}
		case "value":
			err = client.PublishValue(path, os.Args[4])
			if err != nil {
				fmt.Println("Failed to publish value:", err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unknown publish type. Use 'file', 'directory', 'DBTable', 'event', or 'value'.")
			os.Exit(1)
		}
	case "consume":
		var path string
		if len(os.Args) < 2 {
			path = "/" // Default to root if no path is provided
		} else {
			path = os.Args[2]
		}
		data, err := client.Get(path)
		if err != nil {
			fmt.Println("Failed to consume value:", err)
			os.Exit(1)
		}

		if data == nil {
			fmt.Println("No data found at path:", path)
			os.Exit(1)
		}
		switch v := data.(type) {
		case *pb.StringValue:
			fmt.Printf("\"%s\"\n", v.Value)
		case *pb.IntValue:
			fmt.Printf("%d\n", v.Value)
		case *pb.FloatValue:
			fmt.Printf("%.2f\n", v.Value)
		case *pb.DatabaseTable:
			fmt.Printf("Consumed dataset: %v\n", v.TableName)
			fmt.Println(nc.QueryTable(v))
		case *pb.IndividualFile:
			fmt.Printf("Consumed dataset: %v\n", v.FilePath)
			fileStr, err := nc.ReadFile(v.FilePath)
			if err != nil {
				fmt.Println("Failed to read file:", err)
				os.Exit(1)
			}
			fmt.Println(fileStr)
		case *pb.Directory:
			fmt.Printf("Consumed dataset: %v\n", v.DirectoryPath)
		case *pb.EventStream:
			fmt.Printf("Consumed dataset: %v\n", v.Topic)
			messageChan, err := nc.GetEventStream(v)
			if err != nil {
				fmt.Println("Failed to get event stream:", err)
				os.Exit(1)
			}
			for message := range messageChan {
				fmt.Println(string(message))
			}
		default:
			//fmt.Printf("Unknown data type %s", reflect.TypeOf(v).Elem().Name())
			os.Exit(1)
		}
	case "list":
		var path string
		if len(os.Args) > 2 {
			path = os.Args[2]
		} else {
			path = "/" // Default to root if no path is provided
		}
		children, err := client.GetChildren(path)
		if err != nil {
			fmt.Println("Failed to get children:", err)
			os.Exit(1)
		}
		fmt.Printf("Children of path '%s':\n", path)
		for _, child := range children {
			fmt.Println(path + "/" + child)
		}
	default:
		fmt.Println("Unknown command. Use 'publish' or 'consume'.")
		os.Exit(1)
	}
}
