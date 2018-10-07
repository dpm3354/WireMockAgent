package main

import (
	json "encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var watcher *fsnotify.Watcher

type ProcessConfiguration struct {
	Alias     string   `json:"alias"`
	Monitor   string   `json:"monitor"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
}

func spawnProcess(name string, arg ...string) *os.Process {
	log.Printf("Starting \"%v\" with following args: \"%v\"", name, arg)
	cmd := exec.Command(name, arg...)
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

func getProcessConfiguration() *ProcessConfiguration {

	var config string

	if len(os.Args) == 1 {
		config = "config.json"
	} else {
		config = os.Args[1]
	}

	jsonFile, err := os.Open(config)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()


	var managedProcess = ProcessConfiguration{}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &managedProcess)

	if err != nil {
		fmt.Println(err)
	}
	return &managedProcess
}

// see https://medium.com/@skdomino/watch-this-file-watching-in-go-5b5a247cf71f
func main() {

	processConfiguration := *getProcessConfiguration()


	var currentProcess = spawnProcess(processConfiguration.Command, processConfiguration.Arguments...)

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(processConfiguration.Monitor, watchDir); err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Printf("monitoring %v", processConfiguration.Monitor)
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
				currentProcess = spawnProcess(processConfiguration.Command, processConfiguration.Arguments...)
				fmt.Printf("restarted \"%V\"", processConfiguration.Alias)
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
