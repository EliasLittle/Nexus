package main

import (
	"fmt"
	"os"
	"strings"

	nc "nexus/pkg/client"
	pb "nexus/pkg/proto"
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
			fmt.Println("Expected 'dataset', 'event', or 'value' as argument for publish, followed by a path")
			os.Exit(1)
		}
		path := os.Args[3]
		switch os.Args[2] {
		case "dataset":
			dataset := nc.CreateDataset(os.Args[4])
			err = client.PublishDataset(path, dataset)
			if err != nil {
				fmt.Println("Failed to publish dataset:", err)
				os.Exit(1)
			}
		case "event":
			eventStream := nc.CreateEventStream(os.Args[4])
			err = client.PublishEventStream(path, eventStream)
			if err != nil {
				fmt.Println("Failed to publish event stream:", err)
				os.Exit(1)
			}
		case "value":
			value := nc.CreateValue(os.Args[4])
			err = client.PublishValue(path, value)
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
			resultPtr, err := client.GetValue(path)
			if err != nil {
				fmt.Println("Failed to consume value:", err)
				os.Exit(1)
			}

			if resultPtr == nil {
				fmt.Println("No value found at path:", path)
				os.Exit(1)
			}
			switch v := resultPtr.Value.(type) {
			case *pb.Value_StringValue:
				fmt.Printf("\"%s\"\n", v.StringValue.Value)
			case *pb.Value_IntValue:
				fmt.Printf("%d\n", v.IntValue.Value)
			case *pb.Value_FloatValue:
				fmt.Printf("%.2f\n", v.FloatValue.Value)
			default:
				fmt.Println("Unknown value type.")
				os.Exit(1)
			}
		case "dataset":
			fmt.Println("Consuming dataset...")
			dataset, err := client.GetDataset(path)
			if err != nil {
				fmt.Println("Failed to consume dataset:", err)
				os.Exit(1)
			}
			fmt.Printf("Consumed dataset: %v\n", dataset)
		case "event":
			eventStream, err := client.GetEventStream(path)
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
