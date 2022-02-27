package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func CreateFileOnDirectory(text string, fileName string, path ...string) {
	pathResult := filepath.Join(path...)
	path = append(path, fileName)
	pathResultWithFileName := filepath.Join(path...)
	CreateIfNotExistDirectory(pathResult)
	fmt.Println(path)
	file, err := os.Create(pathResultWithFileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created successfully")
	defer file.Close()

	_, err2 := file.WriteString(text)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

func CreateIfNotExistDirectory(directory string) {
	if err := ensureDir(directory); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}
	fmt.Printf("directory created")
}

func ensureDir(dirName string) error {
	path, err := os.Getwd()
	if err != nil {
		log.Println("error msg", err)
	}

	//Create output path
	outPath := filepath.Join(path, dirName)

	//Create dir output using above code
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		os.Mkdir(outPath, 0755)
	}
	return nil
}
