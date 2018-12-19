package configor

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Load(configPath string) (kvMap *map[string]interface{}) {
	newMap := make(map[string]interface{})

	go addConfigPath(configPath, &newMap)

	kvMap = &newMap
	log.Println("Path: ", configPath)
	log.Println("Initial configor...")
	return
}

func addConfigPath(configPath string, kvMap *map[string]interface{}) (err error) {
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

	*kvMap, err = loadConfig(absPath)
	if err != nil {
		return
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
					log.Println("return")
					return
				}
				if event.Name == absPath {
					switch event.Op {
					case fsnotify.Create:
						log.Println(event.Op)
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							return
						}
					case fsnotify.Remove:
						log.Println(event.Op)
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							return
						}
					case fsnotify.Write:
						log.Println(event.Op)
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							return
						}
					case fsnotify.Chmod:
						log.Println(event.Op)
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							return
						}
					case fsnotify.Rename:
						log.Println(event.Op)
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							return
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("return")
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

func loadConfig(path string) (kvMap map[string]interface{}, err error) {
	bytesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	var key []byte
	var value []byte
	var readKey bool
	var hasRead bool
	kvMap = make(map[string]interface{})
	readKey = true
	for _, v := range bytesFile {
		if readKey {
			if v == '\n' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if len(key) > 0 && v == ' ' {
				hasRead = true
				continue
			}
			if v == '=' {
				readKey = false
				hasRead = false
				continue
			}
			if hasRead && v != ' ' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if v != ' ' && !hasRead {
				key = append(key, v)
				continue
			}
		} else {
			if v == '=' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if len(value) > 0 && v == ' ' {
				hasRead = true
				continue
			}
			if v == '\n' {
				readKey = true
				hasRead = false
				kvMap[string(key)] = string(value)
				key = []byte{}
				value = []byte{}
				continue
			}
			if hasRead && v != ' ' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if v != ' ' && !hasRead {
				value = append(value, v)
				continue
			}
		}
	}
	return
}
