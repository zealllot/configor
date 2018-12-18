package configor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

func Load(configPath string) {
	go addConfigPath(configPath)
	log.Println("Path: ", configPath)
	log.Println("Initial configor...")
}

func addConfigPath(configPath string) (err error) {
	fd, err := os.Stat(configPath)
	if err != nil {
		return
	}
	if fd.IsDir() {
		err = errors.New("Path can't be a dir.")
		return
	}

	var absPath string
	if !path.IsAbs(configPath) {
		absPath, _ = filepath.Abs(configPath)
	} else {
		absPath = configPath
	}

	dirPath := path.Dir(configPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("return")
					return
				}
				if event.Name == absPath {
					switch event.Op {
					case fsnotify.Create:
						log.Println(event.Op)
						loadConfig(absPath)
					case fsnotify.Remove:
						log.Println(event.Op)
						loadConfig(absPath)
					case fsnotify.Write:
						log.Println(event.Op)
						loadConfig(absPath)
					case fsnotify.Chmod:
						log.Println(event.Op)
						loadConfig(absPath)
					case fsnotify.Rename:
						log.Println(event.Op)
						loadConfig(absPath)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("return")
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	select {}
	return
}

func loadConfig(path string) (content string, err error) {
	bytesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	content = string(bytesFile)
	log.Println(`Content: 
` + content)
	return
}
