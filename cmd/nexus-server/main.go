package main

import (
	"fmt"
	"net"
	"os"

	"nexus/pkg/logger"
	pb "nexus/pkg/proto"
	ns "nexus/pkg/server"

	"google.golang.org/grpc"
)

func main() {
	// Initialize logger
	log := logger.GetLogger()

	log.Info("Starting Nexus Server...")

	var fullLoadPath, fullSavePath string

	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Error("Invalid number of arguments", "num_args", len(os.Args))
		log.Fatal("Usage: '%s <load_file_path> <save_file_path>' or '%s <save_file_path>'", os.Args[0], os.Args[0])
	}

	if len(os.Args) == 2 {
		fullLoadPath = os.Args[1]
		fullSavePath = os.Args[1]
	} else if len(os.Args) > 2 {
		fullLoadPath = os.Args[1]
		fullSavePath = os.Args[2]
	} else {
		log.Debug("Num args: %d", len(os.Args))
		log.Fatal("Usage: '%s <load_file_path> <save_file_path>' or '%s <save_file_path>'", os.Args[0], os.Args[0])
	}

	expandedPath, err := os.UserHomeDir() // Get the user's home directory
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	fullLoadPath = fmt.Sprintf("%s/%s", expandedPath, fullLoadPath) // Expand the load path
	fullSavePath = fmt.Sprintf("%s/%s", expandedPath, fullSavePath)

	log.Info("Creating new server", "path", fullLoadPath)
	server, err := ns.NewServer(fullLoadPath)
	if err != nil {
		log.Fatal("Failed to create new server", "error", err)
	}

	// Save index on exit
	defer func() {
		if err := server.SaveIndex(fullSavePath); err != nil {
			log.Error("Failed to save index", "error", err)
		} else {
			log.Info("Saved index", "path", fullSavePath)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", ns.DefaultPort))
	if err != nil {
		log.Fatal("Failed to listen", "error", err)
	}

	s := grpc.NewServer()
	pb.RegisterNexusServiceServer(s, server)

	log.Info("Server listening", "port", ns.DefaultPort)

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve", "error", err)
	}
}
