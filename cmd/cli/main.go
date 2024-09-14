package main

import (
	"fmt"
	"io/ioutil"

	"github.com/scirelli/roms-files-cleanup/internal/pkg/log"
)

func main() {
	config := NewConfig()

	files, err := ioutil.ReadDir("./assets/databases")
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		fmt.Println(i, file.Name())
	}
}
