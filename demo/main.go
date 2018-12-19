package main

import (
	"fmt"
	"os"
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
			fmt.Println(k, "=", v)
		}
		fmt.Println("\nonce\n")
	}

	select {}
}
