package configer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func AddConfigPath(configPath string) (err error) {
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

	parentPath := path.Dir(configPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
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
					case fsnotify.Remove:
					case fsnotify.Write:
					case fsnotify.Chmod:
					case fsnotify.Rename:
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

	err = watcher.Add(parentPath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return
}
