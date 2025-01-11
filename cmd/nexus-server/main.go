package main

import (
	"fmt"
	"log"
	"net"
	pb "nexus/pkg/proto"

	ns "nexus/pkg/server"

	"google.golang.org/grpc"
)

func main() {

	// Initialize a new Trie for storing paths
	var index = ns.NewTrie()
	log.Println("Starting Nexus Server...")

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
