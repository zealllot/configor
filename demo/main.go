package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/zealllot/configor"
)

func main() {
	here, _ := os.Getwd()
	configPath := here + "/test"
	m := configor.Load(configPath)
	t := time.Tick(time.Second * 5)

	for {
		<-t
		for k, v := range *m {
			if value, ok := v.(int); ok {
				fmt.Println(reflect.TypeOf(value).Kind())
				fmt.Println(k, "=", value)
			}
			if value, ok := v.(float64); ok {
				fmt.Println(reflect.TypeOf(value).Kind())
				fmt.Println(k, "=", value)
			}
			if value, ok := v.(string); ok {
				fmt.Println(reflect.TypeOf(value).Kind())
				fmt.Println(k, "=", value)
			}

		}
		fmt.Println("\nonce\n")
	}
}
