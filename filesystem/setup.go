package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

var DriveName string = "Drive"
var ConfigName string = "Config.json"

type Options struct {
	Perm int
}

type Drive struct {
	Location string
	Tree     FileBranch
}

type FileTree struct {
	Root FileBranch
	//MetaData along the way
}

type FileBranch struct { //represents folder
	Parent      *FileBranch
	SubBranches []FileBranch
	Leaves      []FileLeaf
	Name        string
}

type BranchNotFoundError struct {
	message    string
	branchname string
}

func (e *BranchNotFoundError) Error() string {
	return e.message
}

func (e *BranchNotFoundError) Handler(drive Drive) {
	filepath.WalkDir(drive.Location, func(path string, info fs.DirEntry, err error) error {
		// filepath
		if path == e.branchname {
			//This
			return nil
		}
		return nil
	})
}

func (b *FileBranch) FindBranchDescending(branchname string) (*FileBranch, error) {
	if b.Name == branchname {
		return b, nil
	}
	for _, subBranch := range b.SubBranches {
		return subBranch.FindBranchDescending(branchname)
	}
	return nil, &BranchNotFoundError{
		message:    "Branch not found in Tree",
		branchname: branchname,
	}

}

func (b *FileBranch) BuildPath() string {
	var pt *FileBranch
	pt = b

	var segments []string
	for pt.Parent != nil {
		segments = append([]string{pt.Name}, segments...)
		pt = pt.Parent
	}

	var fullpath string
	for _, segment := range segments {
		fullpath += "/" + segment
	}

	return fullpath
}

func (b *FileBranch) VerifyPath(path string) bool {
	p1 := filepath.Clean(b.BuildPath())
	p2 := filepath.Clean(path)
	return p1 == p2

}

type PathError struct {
	Message string
	Code    int
}

func (e *PathError) Error() string {
	return e.Message
}

func (b *FileBranch) VerifyPathErr(path string) error {
	if !b.VerifyPath(path) {
		message := path + ", Does not match: " + b.BuildPath()
		err := &PathError{
			Message: message,
		}
		return err
	}
	return nil
}

func (b *FileBranch) AddLeaf(leaf *FileLeaf) error {
	leaf.Branch = b
	b.Leaves = append(b.Leaves, *leaf)
	return nil
}

func (b *FileBranch) AddBranch(branch *FileBranch) error {
	b.Parent = b
	b.SubBranches = append(b.SubBranches, *branch)
	return nil
}

type JSONBranch struct {
	ParentFolder string
	Folders      []string
	File         []FileLeaf
	Name         string
}

type JSONLeaf struct {
	ParentFolder string
	Data         os.FileInfo
	Identifier   string
	Type         string
}

func (b *FileBranch) ToJSON() []byte {
	marshallFriendly := JSONBranch{
		ParentFolder: b.Parent.Name,
		Folders: b.SubBranches.Map(func(ind int, fb FileBranch) {
			return fb.Name
		}),
	}
	json.Marshal(marshallFriendly)
	// type FileBranch struct { //represents folder
	// 	Parent      *FileBranch
	// 	SubBranches []FileBranch
	// 	Leaves      []FileLeaf
	// 	Name        string
	// }

}

type array[T any] []T

func (n array[T]) Map(callback func(int, interface{}) interface{}) []any {
	var mappedArray []interface{}
	for index, element := range n {
		mappedArray = append(mappedArray, callback(index, element))
	}
	return mappedArray
}

type FileLeaf struct { //represents files
	Branch     *FileBranch
	Data       os.FileInfo
	Identifier string
	Type       string
}

func LeafFromFile(file *os.File) (FileLeaf, error) {
	var emptyLeaf FileLeaf // figure out better solution
	stat, err := file.Stat()
	if err != nil {
		return emptyLeaf, err
	}
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return emptyLeaf, err
	}

	return FileLeaf{
		Data:       stat,
		Identifier: HashString(file.Name()),
		Type:       mime.String(),
	}, nil
}

type JSONNode struct {
	Type     string
	Dir      bool
	Data     os.FileInfo
	Children []JSONNode `json:"Children,omitempty"`
}

// func (t *FileTree) JSONCascade() JSONNode {
// 	var jsonnode = JSONNode{
// 		Type: t.Root.Type,
// 		Dir:  t.Root.Dir,
// 		Data: t.Root.Data,
// 	}

// 	if len(t.Root.Children) <= 0 {
// 		//basecase
// 		// JSON, err :=  json.MarshalIndent(jsonnode, "", "\t")
// 		return jsonnode
// 	}
// 	var children []JSONNode
// 	func() {
// 		for _, entry := range t.Root.Children {
// 			root := entry.AsRoot()
// 			children = append(children, root.JSONCascade())
// 		}
// 	}()
// 	jsonnode.Children = children
// 	return jsonnode
// }

// type FSParent struct {
// 	Node     FSNode
// 	Children []FSNode
// }

// func (n *FSParent) Add(node FSNode) {
// 	n.Children = append(n.Children, node)
// }

// func (n *FSParent) BuildPath() string {
// 	var pathSegments []string
// 	pt := *n
// 	for pt.ParentDir != nil {
// 		pathSegments = append([]string{pt.Data.Name()}, pathSegments...)
// 		pt = *pt.ParentDir
// 	}
// 	var fullpath string
// 	for _, segment := range pathSegments {
// 		fullpath += "/" + segment
// 	}
// 	return fullpath
// }

// func (n *FSParent) AsRoot() FileTree {
// 	return FileTree{
// 		Root: *n,
// 	}
// }
// func (n *FSNode) BuildPath() string {
// 	return n.ParentDir.BuildPath() + "/" + n.Data.Name()
// }

// func DirCascade(sourceDir string) NodeInterface {
// 	fileInfo, err := os.Stat(sourceDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	dirEntries, err := os.ReadDir(sourceDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var children []NodeInterface

// 	for _, entry := range dirEntries {
// 		var child []NodeInterface
// 		if entry.IsDir() {
// 			child := DirCascade(filePath)
// 		} else {
// 			child := Node
// 		}

// 	}
// 	root := FSParent{
// 		Data:     fileInfo,
// 		Children: children,
// 	}
// 	return root
// 	// node := createNode(fileInfo)
// }

// func NodeCascade(parent FSParent, entry os.DirEntry) NodeInterface {
// 	// var
// 	filePath := filepath.Join(parent.BuildPath(), entry.Name())
// 	if entry.IsDir() {
// 	}
// 	return FSNode{}
// }

// type NodeInterface interface {
// 	BuildPath() string
// }

func (d *Drive) ToConfig() map[string]any {
	var config map[string]any = map[string]any{
		"Location": d.Location,
	}
	return config
}

func (d *Drive) List(subDirectory string) ([]os.DirEntry, error) {
	var err error
	var dirEntries []os.DirEntry
	fullpath := filepath.Join(d.Location, subDirectory)
	dirEntries, err = os.ReadDir(fullpath)
	return dirEntries, err
}
func (d *Drive) Exists(subDirectory string) bool {
	objectPath := filepath.Join(d.Location, subDirectory)
	_, err := os.Stat(objectPath)
	return os.IsNotExist(err)
}

func (d *Drive) UploadFile(name string, branchName string, file multipart.File) (FileBranch, error) {
	tree := d.Tree
	targetBranch, err := tree.FindBranchDescending(branchName)
	if err != nil {
		//probably should add errorhandler here
		return *targetBranch, err
	}
	osPath := targetBranch.BuildPath()

	newFile, err := os.Create(filepath.Join(osPath, name))
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return *targetBranch, err
	}

	// branch.Update(), but for now:
	// drive.ReconcileTree() idea
	// tree.Reconcile() idea
	leaf, err := LeafFromFile(newFile)
	err = targetBranch.AddLeaf(&leaf)
	if err != nil {
		return *targetBranch, err
	}
	return *targetBranch, nil

	// drive.List(path)

}

func CreateDrive(sourcedir string, options Options) Drive {
	fileInfo, err := os.Stat(sourcedir)
	//probably should scour for an applications folder instead of just writing it in home directory

	fmt.Println(fileInfo)

	if err != nil {
		log.Fatal("Failure accessing directory: ", err)
	}

	drivePath := filepath.Join(sourcedir, DriveName)
	err = os.Mkdir(drivePath, fs.FileMode(options.Perm))
	defer func() {
		if err != nil {
			os.Remove(sourcedir)
			log.Fatal(err)

		}
	}()
	var drive Drive = Drive{
		Location: drivePath,
	}

	return drive
}

// func createNode(data os.FileInfo, parent FSParent) FSNode {
// 	// type FSNode struct {
// 	// 	Parent     *FSNode
// 	// 	Identifier string
// 	// 	Type       string
// 	// 	Dir        bool
// 	// 	Data       os.FileInfo
// 	// 	Children   []FSNode
// 	// }
// 	file, err := os.Open(parent, data.Name())
// 	defer file.Close()

// 	filePath := filepath.Join(parent)
// 	mimetype.DetectFile()

// 	return FSNode{
// 		Type:       data.Mode(),
// 		Dir:        data.IsDir(),
// 		Identifier: HashString(data.Name()),
// 		// Type:
// 		// Data : data,
// 	}

// }

// func createNodeRecursive(data fs.FileInfo, parent FSNode) FSNode {
// 	if data.IsDir() {
// 		return createNodeRecursive(data, parent)
// 	}
// 	return FSNode{
// 		// Data
// 	}
// }

func CheckConfig(sourcedir string) bool {
	configPath := filepath.Join(sourcedir, DriveName, ConfigName)
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}

func WriteConfig(drive Drive, permission int) error {
	config := drive.ToConfig()
	bytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(drive.Location, bytes, fs.FileMode(permission))
	return err
	//os.IsPermission(err)
}

func ReadConfig(drive Drive) {

}
