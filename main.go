package main

import (
	"fmt"
	"log"
	"my-go-project/filesystem"
	"my-go-project/network"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
)

func TestDir() (string, func()) {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirSource := filepath.Join(currentUser.HomeDir, "DriveApp")
	err = os.Mkdir(dirSource, 0700)
	if err != nil {
		log.Fatal(err)
	}
	return dirSource, func() {
		os.RemoveAll(dirSource)
	}
}

func main() {
	dirSource, cleanUp := TestDir()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		fmt.Println(sig)
		cleanUp()
		os.Exit(1)

	}()
	options := filesystem.Options{
		Perm: 0700,
	}
	drive := filesystem.CreateDrive(dirSource, options)
	filesystem.WriteConfig(drive, options.Perm)
	network.Server(drive)
}
