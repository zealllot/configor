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

func addConfigPath(configPath string, kvMap *map[string]interface{}) {
	backupMap := make(map[string]interface{})

	fd, err := os.Stat(configPath)
	if err != nil {
		log.Fatal(err)
	}
	if fd.IsDir() {
		err = errors.New("Path can't be a dir.")
		log.Fatal(err)
	}

	var absPath string
	if !path.IsAbs(configPath) {
		absPath, _ = filepath.Abs(configPath)
	} else {
		absPath = configPath
	}

	backupMap = *kvMap
	*kvMap, err = loadConfig(absPath)
	if err != nil {
		log.Fatal(err)
		*kvMap = backupMap
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
			case event := <-watcher.Events:
				if event.Name == absPath {
					switch event.Op {
					case fsnotify.Create:
						log.Println(event.Op)
						backupMap = *kvMap
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							log.Println(err)
							*kvMap = backupMap
							log.Println("Rollback to the latest version!")
						}
					case fsnotify.Remove:
						log.Println(event.Op)
						backupMap = *kvMap
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							log.Println(err)
							*kvMap = backupMap
							log.Println("Rollback to the latest version!")
						}
					case fsnotify.Write:
						log.Println(event.Op)
						backupMap = *kvMap
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							log.Println(err)
							*kvMap = backupMap
							log.Println("Rollback to the latest version!")
						}
					case fsnotify.Chmod:
						log.Println(event.Op)
						backupMap = *kvMap
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							log.Println(err)
							*kvMap = backupMap
							log.Println("Rollback to the latest version!")
						}
					case fsnotify.Rename:
						log.Println(event.Op)
						backupMap = *kvMap
						*kvMap, err = loadConfig(absPath)
						if err != nil {
							log.Println(err)
							*kvMap = backupMap
							log.Println("Rollback to the latest version!")
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println(err)
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
	var key []rune
	var value []rune
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
				key = append(key, rune(v))
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

				var notInt bool
				var notFloat bool

				for k, vv := range value {
					if !((48 < vv && vv < 57) || vv == 46) {
						notFloat = true
					} else if (k == 0 || k == len(value)-1) && value[k] == 46 {
						notFloat = true
					}
				}

				for _, vv := range value {
					if !(48 < vv && vv < 57) {
						notInt = true
					}
				}

				if !notInt {
					//int
					kvMap[string(key)] = runeToInt(value)
				} else if !notFloat {
					//float
				} else {
					//string
					kvMap[string(key)] = string(value)
				}

				key = []rune{}
				value = []rune{}
				continue
			}
			if hasRead && v != ' ' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if v != ' ' && !hasRead {
				value = append(value, rune(v))
				continue
			}
		}
	}
	return
}
