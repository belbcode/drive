package main

import (
	"fmt"
	"my-go-project/filesystem"
	"os"
)

func main() {
	var err error
	var wd string
	var dirTree filesystem.Tree
	wd, err = os.UserHomeDir()

	if err != nil {
		fmt.Println(err)
	}
	rootNode, err := filesystem.RecursiveBuildFromRoot(wd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rootNode)
	dirTree = filesystem.Tree{Root: rootNode}
	treejson := dirTree.ToJSON()

	os.WriteFile("tree.json", treejson, 0700)

	// results := dirTree.Search("hello")
	// for _, node := range results {
	// 	fmt.Println(node.Info.Name())
	// }
	// fmt.Println(*dirStruct.SubDirectories)

	// // reader := dirStruct
	// pointer := *dirStruct.SubDirectories

	// for {
	// 	fmt.Scanln()
	// 	if !(len(pointer) > 0) {
	// 		fmt.Println("End of dir")
	// 		break
	// 	}
	// 	pointer = *pointer[0].SubDirectories
	// 	fmt.Println(pointer)
	// }
}
