package main

import (
	"fmt"
	"log"
	"net"
	pb "nexus/pkg/proto"
	"os"
	"os/signal"

	ns "nexus/pkg/server"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Nexus Server...")

	var loadFilePath, saveFilePath string
	if len(os.Args) == 2 {
		loadFilePath = os.Args[1]
		saveFilePath = os.Args[1]
	} else if len(os.Args) > 2 {
		loadFilePath = os.Args[1]
		saveFilePath = os.Args[2]
	} else {
		log.Printf("Num args: %d", len(os.Args))
		log.Fatalf("Usage: '%s <load_file_path> <save_file_path>' or '%s <save_file_path>'", os.Args[0], os.Args[0])
	}

	// Generate full paths for load and save file
	expandedPath, err := os.UserHomeDir() // Get the user's home directory
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	fullLoadPath := fmt.Sprintf("%s/%s", expandedPath, loadFilePath) // Expand the load path
	fullSavePath := fmt.Sprintf("%s/%s", expandedPath, saveFilePath) // Expand the save path

	// Initialize a new Trie for storing paths
	var index *ns.Trie
	if loadFilePath != "" {
		log.Printf("Loading index from %s", fullLoadPath)
		index, err = ns.NewTrie(fullLoadPath) // Load from the provided file path
	} else {
		index, err = ns.NewTrie() // No load path provided
	}

	if err != nil {
		log.Fatalf("Failed to create new Trie: %v", err)
	}

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Defer saving the trie to the save file on shutdown
	go func() {
		<-signalChan            // Wait for interrupt signal
		if saveFilePath != "" { // Only save if a save path was provided
			if err := index.SaveToDisk(fullSavePath); err != nil {
				log.Printf("Failed to save index: %v", err)
			} else {
				log.Printf("Saved index to %s", fullSavePath)
			}
		}
		os.Exit(0) // Exit gracefully
	}()

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", ns.DefaultPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterNexusServiceServer(s, &ns.NexusServer{Index: index})

	log.Printf("Server listening at :%d", ns.DefaultPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
