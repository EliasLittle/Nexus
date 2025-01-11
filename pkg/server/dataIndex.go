package server

import (
	"fmt"
	"log"
	pb "nexus/pkg/proto"
	"strings"
)

type TrieNode struct {
	Children    map[string]*TrieNode
	IsEndOfPath bool
	Value       interface{} // Store different types of values
}

type Trie struct {
	Root *TrieNode
}

// NewTrie initializes a new Trie
func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[string]*TrieNode)}}
}

// Insert adds a new path to the Trie with an associated value
func (t *Trie) Insert(path string, value interface{}) {
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			node.Children[segment] = &TrieNode{Children: make(map[string]*TrieNode)}
		}
		node = node.Children[segment]
	}
	// TODO: Think about if we want to store data in intermediate nodes and not just at the end of the path
	node.IsEndOfPath = true
	node.Value = value // Store the value at the end of the path
}

// Search checks if a path exists in the Trie
func (t *Trie) Search(path string) bool {
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			return false
		}
		node = node.Children[segment]
	}
	return node.IsEndOfPath
}

// Traverse prints all paths in the Trie
func (t *Trie) Traverse() {
	t.traverseHelper(t.Root, "")
}

func (t *Trie) traverseHelper(node *TrieNode, prefix string) {
	if node.IsEndOfPath {
		fmt.Println(prefix, " : ", node.Value)
	}
	for segment, child := range node.Children {
		t.traverseHelper(child, prefix+"/"+segment)
	}
}

// splitPath is a helper function to split the path into segments
func splitPath(path string) []string {
	// This can be customized for different segmenting logic
	return strings.Split(strings.Trim(path, "/"), "/")
}

// GetType returns the type of data stored at the given path
func (t *Trie) GetType(path string) (string, error) {
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			return "", fmt.Errorf("path not found: %s", path)
		}
		node = node.Children[segment]
	}
	if node.IsEndOfPath {
		switch node.Value.(type) {
		case *pb.EventStream:
			return "EventStream", nil
		case *pb.Dataset:
			return "Dataset", nil
		case *pb.DirectValue:
			return "Value", nil
		default:
			return "Unknown", nil
		}
	}
	return "", fmt.Errorf("path does not point to a data type: %s", path)
}

// GetChildren retrieves all direct children of a given path
func (t *Trie) GetChildren(path string) []string {
	if path == "/" { // Check if the path is the root
		log.Printf("Children of root path: %v\n", t.Root.Children)
		var children []string
		for childSegment := range t.Root.Children {
			children = append(children, childSegment)
		}
		log.Printf("Children of path '%s': %v\n", path, children)
		return children
	}

	node := t.Root
	segments := splitPath(path)
	for _, segment := range segments {
		if child, exists := node.Children[segment]; exists {
			node = child
		} else {
			log.Printf("No children found for path: %s\n", path)
			return []string{} // No children found
		}
	}

	// Collect all child paths
	var children []string
	for childSegment := range node.Children {
		children = append(children, childSegment)
	}

	log.Printf("Children of path '%s': %v\n", path, children)
	return children
}

// Get retrieves the value stored at the given path in the Trie
func (t *Trie) Get(path string) (interface{}, error) {
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			return nil, fmt.Errorf("path does not exist: %s", path) // Return an error if the path doesn't exist
		}
		node = node.Children[segment]
	}
	if node.IsEndOfPath {
		return node.Value, nil // Return the value if it exists
	}
	return nil, fmt.Errorf("path exists but has no value: %s", path) // Return an error if the path exists but has no value
}

/*
func dataIndexExample() {
	// Connect to the Nexus server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Example usage of GetType
	dataType, err := pathTrie.GetType("/events/sensors")
	if err != nil {
		log.Fatalf("Error getting type: %v", err)
	}
	log.Printf("Data type at path '/events/sensors': %s", dataType)
}
*/
