package server

import (
	"encoding/json"
	"fmt"
	"log"
	pb "nexus/pkg/proto"
	"os"
	"reflect"
	"strings"
)

type TrieNode struct {
	Children    map[string]*TrieNode
	IsEndOfPath bool
	Value       interface{} // Store different types of values
	ValueType   string      // Type identifier for the value
}

type Trie struct {
	Root *TrieNode
}

// NewTrie initializes a new Trie, optionally loading from a file
func NewTrie(filepath ...string) (*Trie, error) {
	if len(filepath) > 0 {
		log.Printf("Loading from disk: %s", filepath[0])
		trie, err := LoadFromDisk(filepath[0]) // Load from the provided file path
		if err != nil {
			log.Printf("Error loading from disk: %v", err)
			return nil, err
		}
		trie.Traverse()
		return trie, nil
	} else {
		trie := &Trie{Root: &TrieNode{Children: make(map[string]*TrieNode)}}
		return trie, nil
	}
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
	node.ValueType = reflect.TypeOf(value).Elem().Name()

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
		case *pb.StringValue:
			return "StringValue", nil
		case *pb.IntValue:
			return "IntValue", nil
		case *pb.FloatValue:
			return "FloatValue", nil
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
	t.Traverse()
	log.Printf("Get| Path: %s", path)
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			return nil, fmt.Errorf("path does not exist: %s", path) // Return an error if the path doesn't exist
		}
		node = node.Children[segment]
	}
	log.Printf("Get| Node: %v", node)
	if node.IsEndOfPath {
		log.Printf("Get| Is end of path")
		return node.Value, nil // Return the value if it exists
	}
	log.Printf("Get| Is not end of path")
	return nil, fmt.Errorf("path exists but has no value: %s", path) // Return an error if the path exists but has no value
}

// SaveToDisk saves the current Trie to a file
func (t *Trie) SaveToDisk(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(t)
}

// LoadFromDisk initializes a Trie from a file
func LoadFromDisk(filename string) (*Trie, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var trie Trie
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&trie); err != nil {
		log.Printf("Error decoding Trie: %v", err)
		return nil, err
	}

	log.Printf("Traversing before loading values")
	trie.Traverse()

	// Load values based on their type identifiers
	log.Printf("Loading values")
	loadValues(&trie)

	log.Printf("Traversing after loading values")
	trie.Traverse()

	return &trie, nil
}

// loadValues reconstructs the values based on their type identifiers
func loadValues(trie *Trie) {
	loadNodeValues(trie.Root)
	log.Printf("Internal traversal")
	trie.Traverse()
}

func loadNodeValues(node *TrieNode) {
	if node.IsEndOfPath {
		log.Printf("Loading values for node: %s", node.ValueType)
		switch node.ValueType {
		case "EventStream":
			node.Value = &pb.EventStream{}
		case "Dataset":
			// Extract the actual dataset value from the representation
			log.Printf("Found dataset")
			switch node.Value.(type) {
			case string:
				node.Value = parseDatasetValue(node.Value.(string))
			case map[string]interface{}:
				log.Printf("Found dataset as map")
				//node.Value = parseDatasetValue(node.Value.(map[string]interface{}))
			default:
				log.Printf("Failed to extract dataset value: %v", node.Value)
			}
		case "StringValue":
			log.Printf("Found string value")
			// Extract the actual string value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Printf("Successfully found string value as map")
				node.Value = &pb.StringValue{Value: valueMap["value"].(string)}
			} else {
				log.Printf("Failed to extract string value: %v", node.Value)
			}
		case "IntValue":
			log.Printf("Found int value")
			// Extract the actual integer value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Printf("Successfully found int value as map")
				node.Value = &pb.IntValue{Value: int32(valueMap["value"].(int))}
			}
		case "FloatValue":
			log.Printf("Found float value")
			// Extract the actual float value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Printf("Successfully found float value as map")
				node.Value = &pb.FloatValue{Value: float32(valueMap["value"].(float64))}
			}
		default:
			log.Printf("Unknown type: %s", node.ValueType)
			node.Value = nil // Handle unknown types
		}
	} else {
		log.Printf("Loading values for node: %s", node.ValueType)
	}
	for _, child := range node.Children {
		loadNodeValues(child)
	}
}

// Helper function to parse the dataset value from the string representation
func parseDatasetValue(str string) *pb.Dataset {
	// Assuming the format is "map[Dataset:map[IndividualFile:map[column_names:[id name email] file_path:./tests/example_a.csv file_type:csv]]]"
	parts := strings.Split(str, "map[IndividualFile:")
	if len(parts) < 2 {
		return nil
	}
	//individualFilePart := strings.Trim(parts[1], "]")

	// Here you would parse individualFilePart to extract the actual values
	// For demonstration, we will just return a new Dataset with dummy values
	return &pb.Dataset{
		Dataset: &pb.Dataset_IndividualFile{
			IndividualFile: &pb.IndividualFile{
				ColumnNames: []string{"id", "name", "email"},
				FilePath:    "./tests/example_a.csv",
				FileType:    "csv",
			},
		},
	}
}
