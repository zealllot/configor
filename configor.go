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
	for idBytes, aByte := range bytesFile {
		if readKey {
			if aByte == '\n' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if len(key) > 0 && aByte == ' ' {
				hasRead = true
				continue
			}
			if aByte == '=' {
				readKey = false
				hasRead = false
				continue
			}
			if hasRead && aByte != ' ' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if aByte != ' ' && !hasRead {
				key = append(key, rune(aByte))
				continue
			}
		} else {
			if aByte == '=' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if len(value) > 0 && aByte == ' ' {
				hasRead = true
				continue
			}
			if aByte == '\n' || idBytes == len(bytesFile)-1 {
				if idBytes == len(bytesFile)-1 && aByte != ' ' && !hasRead {
					value = append(value, rune(aByte))
				}
				readKey = true
				hasRead = false

				var notInt bool
				var notFloat bool

				for idValue, aRune := range value {
					if !((48 < aRune && aRune < 57) || aRune == 46) {
						notFloat = true
					} else if (idValue == 0 || idValue == len(value)-1) && value[idValue] == 46 {
						notFloat = true
					}
				}

				for _, aRune := range value {
					if !(48 < aRune && aRune < 57) {
						notInt = true
					}
				}

				if !notInt {
					//int
					kvMap[string(key)] = runesToInt(value)
				} else if !notFloat {
					//float
					kvMap[string(key)] = runesToFloat64(value)
				} else {
					//string
					if _, ok := kvMap[string(value)]; ok {
						kvMap[string(key)] = kvMap[string(value)]
					} else {
						kvMap[string(key)] = string(value)
					}
				}

				key = key[:0]
				value = value[:0]
				continue
			}
			if hasRead && aByte != ' ' {
				err = errors.New("Incorrect config file! ")
				return
			}
			if aByte != ' ' && !hasRead {
				value = append(value, rune(aByte))
				continue
			}
		}
	}
	return
}
