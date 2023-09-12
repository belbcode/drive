package filesystem

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var DriveName string = "Drive"
var ConfigName string = "Config.json"

type Options struct {
	Perm int
}

type Drive struct {
	Location string
	Tree     FileTree
}

type FileTree struct {
	Root FSNode //where Dir : true
}

type JSONNode struct {
	Type     string
	Dir      bool
	Data     os.FileInfo
	Children []JSONNode `json:"Children,omitempty"`
}

func (t *FileTree) JSONCascade() JSONNode {
	var jsonnode = JSONNode{
		Type: t.Root.Type,
		Dir:  t.Root.Dir,
		Data: t.Root.Data,
	}

	if len(t.Root.Children) <= 0 {
		//basecase
		// JSON, err :=  json.MarshalIndent(jsonnode, "", "\t")
		return jsonnode
	}
	var children []JSONNode
	func() {
		for _, entry := range t.Root.Children {
			root := entry.AsRoot()
			children = append(children, root.JSONCascade())
		}
	}()
	jsonnode.Children = children
	return jsonnode
}

type FileType int

const (
	Photo FileType = iota
	Media
	Document
	File
	Dir
)

type FSNode struct {
	Parent   *FSNode
	Type     string
	Dir      bool
	Data     os.FileInfo
	Children []FSNode
}

func (n *FSNode) Add(node FSNode) {
	n.Children = append(n.Children, node)
}

func (n *FSNode) AsRoot() FileTree {
	return FileTree{
		Root: *n,
	}
}

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
