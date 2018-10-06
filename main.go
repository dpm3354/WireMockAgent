package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var watcher *fsnotify.Watcher

func spawnProcess() *os.Process {
	cmd := exec.Command("java", "-jar", "wiremock-standalone-2.19.0.jar", "--port", "8082", "--global-response-templating")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return cmd.Process
}

func kill(process *os.Process) {
	if err := process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
	}
}

// see https://medium.com/@skdomino/watch-this-file-watching-in-go-5b5a247cf71f
func main() {

	var currentProcess = spawnProcess()

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	path := os.Args[1]
	if err := filepath.Walk(path, watchDir); err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Printf("monitoring %v", path)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)
				kill(currentProcess)
				currentProcess = spawnProcess()
				fmt.Printf("restarted wiremock")
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
