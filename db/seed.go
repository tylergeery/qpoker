package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	var files []string

	db, err := sql.Open("postgres", os.Getenv("PG_CONNECTION"))
	if err != nil {
		panic(err)
	}

	root := "/go/db/seed/"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Printf("Applying file: %s\n", file)
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(string(contents))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Applied all migration files")
}
