package server

import (
	"encoding/json"
	"fmt"
	"nexus/pkg/logger"
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
	log := logger.GetLogger()
	log.Info("Creating Trie")
	if len(filepath) > 0 {
		log.Info("Loading from disk", "path", filepath[0])
		trie, err := LoadFromDisk(filepath[0]) // Load from the provided file path
		if err != nil {
			log.Error("Error loading from disk", "error", err)
			return nil, err
		}
		trie.Traverse()
		return trie, nil
	} else {
		log.Info("Creating new Trie")
		trie := &Trie{Root: &TrieNode{Children: make(map[string]*TrieNode), ValueType: "InternalNode"}}
		return trie, nil
	}
}

// Insert adds a new path to the Trie with an associated value
func (t *Trie) Insert(path string, value interface{}) {
	log := logger.GetLogger()
	node := t.Root
	segments := splitPath(path) // Customizable segmenter
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			node.Children[segment] = &TrieNode{Children: make(map[string]*TrieNode), ValueType: "InternalNode"}
		}
		node = node.Children[segment]
	}
	// TODO: Think about if we want to store data in intermediate nodes and not just at the end of the path
	node.IsEndOfPath = true
	node.Value = value // Store the value at the end of the path
	valueType := reflect.TypeOf(value).Elem().Name()
	log.Info("Inserted value", "value", value, "type", valueType)
	node.ValueType = valueType

}

// Search checks if a path exists in the Trie
func (t *Trie) Search(path string) bool {
	node := t.Root
	segments := splitPath(path)
	for _, segment := range segments {
		if _, exists := node.Children[segment]; !exists {
			return false
		}
		node = node.Children[segment]
	}
	return node.IsEndOfPath
}

// Traverse prints the entire Trie structure
func (t *Trie) Traverse() {
	log := logger.GetLogger()
	log.Debug("Traversing Trie")
	t.traverseHelper(t.Root, "")
}

func (t *Trie) traverseHelper(node *TrieNode, prefix string) {
	log := logger.GetLogger()
	if node.IsEndOfPath {
		log.Debug("Node", "prefix", prefix, "value", node.Value, "type", node.ValueType)
	}
	for segment, child := range node.Children {
		t.traverseHelper(child, prefix+"/"+segment)
	}
}

// splitPath is a helper function to split the path into segments
func splitPath(path string) []string {
	// Remove leading and trailing slashes
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// GetType returns the type of value stored at a path
func (t *Trie) GetType(path string) (string, error) {
	log := logger.GetLogger()
	node, err := t.GetNode(path)
	if err != nil {
		return "", err
	}
	log.Debug("Children of root path", "children", t.Root.Children)
	return node.ValueType, nil
}

// GetChildren returns a list of child paths for a given path
func (t *Trie) GetChildren(path string) []*pb.ChildInfo {
	log := logger.GetLogger()
	node := t.Root
	if path != "/" {
		segments := splitPath(path)
		for _, segment := range segments {
			if child, exists := node.Children[segment]; exists {
				node = child
			} else {
				log.Debug("No children found for path", "path", path)
				return []*pb.ChildInfo{}
			}
		}
	}

	children := make([]*pb.ChildInfo, 0, len(node.Children))
	for segment, child := range node.Children {
		children = append(children, &pb.ChildInfo{
			Name:        segment,
			Type:        child.ValueType,
			NumChildren: int32(len(child.Children)),
		})
	}
	log.Debug("Children of path", "path", path, "children", children)
	return children
}

// GetNode returns the TrieNode at a given path
func (t *Trie) GetNode(path string) (*TrieNode, error) {
	log := logger.GetLogger()
	log.Debug("Getting node", "path", path)
	node := t.Root
	if path != "/" {
		segments := splitPath(path)
		for _, segment := range segments {
			if child, exists := node.Children[segment]; exists {
				log.Debug("Found child", "child", child)
				node = child
			} else {
				log.Debug("Path not found", "path", path)
				return nil, fmt.Errorf("path not found: %s", path)
			}
		}
	}
	return node, nil
}

// Delete deletes a path from the Trie
func (t *Trie) Delete(path string) (bool, error) {
	log := logger.GetLogger()
	log.Debug("Deleting path", "path", path)
	node := t.Root
	segments := splitPath(path)

	for i, segment := range segments {
		if i == len(segments)-1 {
			// If it's the last segment, check if it exists and delete it
			if _, exists := node.Children[segment]; exists {
				delete(node.Children, segment)
				log.Debug("Deleted segment", "segment", segment)
				return true, nil
			} else {
				log.Debug("Path not found, nothing to delete", "path", path)
				return false, fmt.Errorf("path not found: %s", path)
			}
		}

		// Move to the next segment
		if child, exists := node.Children[segment]; exists {
			node = child
		} else {
			log.Debug("Path not found, nothing to delete", "path", path)
			return false, fmt.Errorf("path not found: %s", path)
		}
	}

	log.Debug("Path deleted", "path", path)
	return false, fmt.Errorf("path not found: %s", path)
}

// SaveToDisk saves the Trie to a file
func (t *Trie) SaveToDisk(filename string) error {
	log := logger.GetLogger()
	file, err := os.Create(filename)
	if err != nil {
		log.Error("Error encoding Trie", "error", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(t)
}

// LoadFromDisk loads a Trie from a file
func LoadFromDisk(filename string) (*Trie, error) {
	log := logger.GetLogger()
	log.Debug("Loading from disk", "path", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var trie Trie
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&trie); err != nil {
		log.Error("Error decoding Trie: %v", err)
		return nil, err
	}

	log.Debug("Traversing before loading values")
	trie.Traverse()

	log.Debug("Loading values")
	loadValues(&trie)

	log.Debug("Traversing after loading values")
	trie.Traverse()

	return &trie, nil
}

func loadValues(trie *Trie) {
	log := logger.GetLogger()
	log.Debug("Internal traversal")
	loadNodeValues(trie.Root)
}

func loadNodeValues(node *TrieNode) {
	log := logger.GetLogger()
	if node.IsEndOfPath {
		log.Info("Loading values for node: %s", node.ValueType)
		switch node.ValueType {
		case "EventStream":
			log.Debug("Found event stream")
			node.Value = parseEventStream(node.Value.(map[string]interface{}))
		case "IndividualFile":
			log.Debug("node.Value: %v", node.Value)
			log.Debug("Found individual file")
			node.Value = parseIndividualFile(node.Value)
		case "Directory":
			log.Debug("Found directory")
			node.Value = parseDirectory(node.Value.(map[string]interface{}))
		case "DatabaseTable":
			log.Debug("Found database table")
			node.Value = parseDatabaseTable(node.Value.(map[string]interface{}))
		case "StringValue":
			log.Debug("Found string value")
			// Extract the actual string value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Debug("Successfully found string value as map")
				node.Value = &pb.StringValue{Value: valueMap["value"].(string)}
			} else {
				log.Error("Failed to extract string value", "value", node.Value)
			}
		case "IntValue":
			log.Debug("Found int value")
			// Extract the actual integer value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Debug("Successfully found int value as map")
				node.Value = &pb.IntValue{Value: int32(valueMap["value"].(int))}
			}
		case "FloatValue":
			log.Debug("Found float value")
			// Extract the actual float value from the representation
			if valueMap, ok := node.Value.(map[string]interface{}); ok {
				log.Debug("Successfully found float value as map")
				node.Value = &pb.FloatValue{Value: float32(valueMap["value"].(float64))}
			}
		default:
			log.Error("Unknown type: %s", node.ValueType)
			node.Value = nil // Handle unknown types
		}
	}

	log.Debug("Loading values for node", "type", node.ValueType)
	for _, child := range node.Children {
		loadNodeValues(child)
	}
}

func parseIndividualFile(nodeValue interface{}) *pb.IndividualFile {
	log := logger.GetLogger()
	if valueMap, ok := nodeValue.(map[string]interface{}); ok {
		log.Debug("Successfully found individual file as map")
		var columnNames []string
		if cols, exists := valueMap["column_names"]; exists {
			colStr := cols.(string)
			columnNames = strings.Split(colStr, " ")
		}
		return &pb.IndividualFile{
			FilePath:    valueMap["file_path"].(string),
			FileType:    valueMap["file_type"].(string),
			ColumnNames: columnNames,
		}
	}
	return nil
}

func parseDirectory(nodeValue interface{}) *pb.Directory {
	log := logger.GetLogger()
	if valueMap, ok := nodeValue.(map[string]interface{}); ok {
		log.Debug("Successfully found directory as map")
		return &pb.Directory{
			FileType:      valueMap["file_type"].(string),
			DirectoryPath: valueMap["directory_path"].(string),
			FileCount:     int32(valueMap["file_count"].(float64)),
		}
	}
	return nil
}

func parseDatabaseTable(nodeValue interface{}) *pb.DatabaseTable {
	log := logger.GetLogger()
	if valueMap, ok := nodeValue.(map[string]interface{}); ok {
		log.Debug("Successfully found database table as map")
		return &pb.DatabaseTable{
			DbType:    valueMap["db_type"].(string),
			Host:      valueMap["host"].(string),
			Port:      int32(valueMap["port"].(float64)),
			DbName:    valueMap["db_name"].(string),
			TableName: valueMap["table_name"].(string),
		}
	}
	return nil
}

func parseEventStream(nodeValue interface{}) *pb.EventStream {
	log := logger.GetLogger()
	if valueMap, ok := nodeValue.(map[string]interface{}); ok {
		log.Debug("Successfully found event stream as map")
		return &pb.EventStream{
			Server: valueMap["server"].(string),
			Topic:  valueMap["topic"].(string),
		}
	}
	return nil
}
