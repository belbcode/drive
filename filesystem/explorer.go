package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type NotDirError struct {
	message string
}

func (e *NotDirError) Error() string {
	return e.message
}

func Explorer() {

}

func getFilePath(fileName string) (string, error) {
	// Check if the file exists.
	_, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}

	// Get the absolute path to the file.
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func Navigate(dir os.DirEntry) ([]os.DirEntry, error) {

	var filepath string
	var dirEntries []os.DirEntry
	var err error

	//check is entry is a directory
	if !dir.IsDir() {
		return nil, &NotDirError{message: "Path: " + dir.Name() + "does not point to a directory."}
	}

	//get full path
	filepath, err = getFilePath(dir.Name())

	if err != nil {
		return nil, err
	}

	//read and return project contents
	dirEntries, err = os.ReadDir(filepath)
	if err != nil {
		return nil, err
	}

	return dirEntries, nil
}

type Directory struct {
	Name           string
	Files          *[]string
	SubDirectories *[]Directory
}

type Node struct {
	Info       os.FileInfo
	Identifier string
	Parent     *Node
	Children   []*Node
	Leaf       bool
}

type NodeJSON struct {
	Info     []byte
	Children []NodeJSON
	Leaf     bool
}

func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) SubTree() Tree {
	return Tree{Root: n}
}

func (n *Node) toJsonStruct() NodeJSON {
	var jsonChildren []NodeJSON
	for _, node := range n.Children {
		jsonChildren = append(jsonChildren, node.toJsonStruct())
	}
	jInfo, err := json.MarshalIndent(n.Info, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	jsonStruct := NodeJSON{
		Info:     jInfo,
		Children: jsonChildren,
		Leaf:     n.Leaf,
	}
	return jsonStruct
}

type Tree struct {
	Root *Node
}

func (t *Tree) Traverse(callback func(Node)) {
	pointer := t.Root
	for _, child := range pointer.Children {
		callback(*child)
		if len(child.Children) > 0 {
			subtree := child.SubTree()
			subtree.Traverse(callback)
		}
	}
}

func (t *Tree) ToJSON() []byte {
	jsonStruct := t.Root.toJsonStruct()
	tJson, err := json.MarshalIndent(jsonStruct, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	return tJson
}

func HashString(input string) string {
	hasher := sha256.New()

	// Write the string bytes to the hash object
	hasher.Write([]byte(input))

	// Get the hashed bytes
	hashedBytes := hasher.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	return hex.EncodeToString(hashedBytes)
}

func LeafFromPath(filepath string) (Node, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return Node{}, err
	}
	var children []*Node
	return Node{
		Info:       fileInfo,
		Identifier: HashString(filepath),
		Parent:     nil,
		Leaf:       true,
		Children:   children,
	}, nil
}

func RecursiveBuildFromRoot(rootDirectory string) (*Node, error) {

	parent, err := LeafFromPath(rootDirectory)
	if err != nil {
		fmt.Println("error", rootDirectory)
		return &parent, nil
	}

	entries, err := os.ReadDir(rootDirectory)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		filepath := filepath.Join(rootDirectory, entry.Name())
		if entry.IsDir() {
			node, err := RecursiveBuildFromRoot(filepath)
			if err != nil {
				fmt.Println(err, "Failed to add: ", filepath)
				continue
			}
			node.Leaf = false
			parent.AddChild(node)
		} else {
			node, err := LeafFromPath(filepath)
			if err != nil {
				fmt.Println(err, "Failed to add: ", filepath)
				continue
			}
			node.Parent = &parent
			parent.AddChild(&node)
		}
	}
	return &parent, nil
}

// func BuildTree(rootDirectory string, parent Node) {
// 	fileInfo, err := os.Stat(rootDirectory)
// 	if err != nil {
// 		return
// 	}
// 	entries, err := os.ReadDir(rootDirectory)
// 	if err != nil {
// 		return
// 	}

// 	var children []*Node

// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			filepath, err := getFilePath(entry.Name())
// 			if err != nil {
// 				return
// 			}

// 		}
// 	}
// 	rootNode := Node{
// 		Info:     fileInfo,
// 		Parent:   nil,
// 		Children: children,
// 	}
// }

func Peek(currentDirectory string) (Directory, error) {

	var subdirectories []Directory
	var files []string
	structure := &Directory{
		Name:           filepath.Base(currentDirectory),
		SubDirectories: &subdirectories,
		Files:          &files,
	}

	var filepath string
	var dirEntries []os.DirEntry
	var err error

	filepath, err = getFilePath(currentDirectory)

	if err != nil {
		return *structure, err
	}

	dirEntries, err = os.ReadDir(filepath)
	if err != nil {
		return *structure, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			directory, err := Peek(entry.Name())

			if err != nil {
				panic(err)
			}
			subdirectories = append(subdirectories, directory)
		} else {
			fmt.Println(entry.Name())

			files = append(files, entry.Name())
		}
	}
	return *structure, nil
}

func Traverse(directory string) {

}
