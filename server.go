package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
)

//go:embed data
var content embed.FS

func main() {
	err := fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Printf("%s\n", path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
}
