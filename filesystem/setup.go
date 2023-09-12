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
	Tree     interface{}
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

func CreateDrive(sourcedir string, options Options) Drive {
	fileInfo, err := os.Stat(sourcedir)

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
	configPath := filepath.Join(drive.Location, DriveName)
	err = os.WriteFile(configPath, bytes, fs.FileMode(permission))
	return err
	//os.IsPermission(err)
}

func ReadConfig(drive Drive) {

}
