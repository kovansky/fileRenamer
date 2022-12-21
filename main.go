package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var (
		err                                     error
		dirInput, extInput, nameInput, numInput string
		info                                    os.FileInfo
		ext, dirPath                            string
		startNumber, renamed                    int
		files                                   []os.FileInfo
		mapping                                 = make(map[string]string)
		mappingFile                             *os.File
		jsonMapping                             []byte
	)

START:
	renamed = 0

	// Header
	fmt.Println("=== Welcome to F4 Bulk File Renamer ===")
	fmt.Println("----------------------------------")

	// Start with directory
	fmt.Println("Please, pass the directory in which you want to rename files")
	fmt.Print("-> ")

	dirInput, err = reader.ReadString('\n')
	dirInput = strings.ReplaceAll(dirInput, "\n", "")
	if err != nil {
		panic(err)
	}

	dirPath = strings.TrimSpace(path.Clean(dirInput))

	info, err = os.Stat(dirPath)

	if os.IsNotExist(err) {
		fmt.Println("Given path does not exist.")
		goto QUIT
	} else if err != nil {
		fmt.Println(dirInput)
		panic(err)
	}

	if !info.IsDir() {
		fmt.Println("Given path is not a directory.")
		goto QUIT
	}

	// Define, which files you want to rename (extension)
	fmt.Println("Please, specify the extension of files you want to rename without dot (for example: jpg)")
	fmt.Print("-> ")

	extInput, err = reader.ReadString('\n')
	extInput = strings.TrimSpace(extInput)
	if err != nil {
		panic(err)
	}

	ext = strings.TrimPrefix(extInput, ".")

	// Now, the new files name
	fmt.Println("Now specify how the new files should be named (it will be used as: {name}{number}.{extension})")
	fmt.Print("-> ")

	nameInput, err = reader.ReadString('\n')
	nameInput = strings.TrimSpace(nameInput)
	if err != nil {
		panic(err)
	}

	// Starting number
	fmt.Println("We will start numbering with...")
	fmt.Print("-> ")

	numInput, err = reader.ReadString('\n')
	numInput = strings.TrimSpace(numInput)
	if err != nil {
		panic(err)
	}

	startNumber, err = strconv.Atoi(numInput)
	if err != nil {
		panic(err)
	}

	// Now, lets start the game
	files, err = ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "."+ext) {
			oldName := path.Join(dirPath, f.Name())
			newName := path.Join(dirPath, nameInput+strconv.Itoa(startNumber)+"."+ext)

			err = os.Rename(oldName, newName)

			mapping[oldName] = newName

			if err != nil {
				fmt.Println(fmt.Sprintf("Error while renaming file: %s", oldName))
				fmt.Println("An error occurred after " + strconv.Itoa(renamed) + " files.")
				panic(err)
			}
			startNumber++
			renamed++
		}
	}

	fmt.Println("Successfully renamed " + strconv.Itoa(renamed) + " files.")

	// Write mapping to json file
	mappingFile, err = os.Create(path.Join(dirPath, "rename-mapping.json"))
	if err != nil {
		fmt.Println("Could not create a rename-mapping.json file")
		panic(err)
	}

	jsonMapping, err = json.Marshal(mapping)
	if err != nil {
		fmt.Println("Could not convert files mapping to json")
		panic(err)
	}

	_, err = mappingFile.Write(jsonMapping)
	if err != nil {
		fmt.Println("Could not write mapping to json file")
		panic(err)
	}

QUIT:
	fmt.Println("Enter q to quit, any other key to start over.")
	fmt.Print("-> ")

	quitInput, err := reader.ReadString('\n')
	quitInput = strings.TrimSpace(quitInput)
	if err != nil {
		panic(err)
	}

	if quitInput == "q" {
		return
	} else {
		goto START
	}
}
