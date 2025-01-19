package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type CheckDir struct {
	Name    string      `json:"dir"`
	Files   []string    `json:"files"`
	Folders []*CheckDir `json:"folders"`
}

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	var setsNumber int

	//file, fErr := os.Open("5")
	//if fErr != nil {
	//	log.Fatalf("file open error: %s", fErr.Error())
	//}

	in = bufio.NewReader(os.Stdin)
	//in = bufio.NewReader(file)
	out = bufio.NewWriter(os.Stdout)
	defer func() {
		err := out.Flush()
		if err != nil {
			panic(err)
		}
	}()

	_, err := fmt.Fscan(in, &setsNumber)
	if err != nil {
		log.Fatalf("Count of sets scan error: %s", err.Error())
	}

	for i := 0; i < setsNumber; i++ {
		var lineCount int
		_, err := fmt.Fscan(in, &lineCount)
		if err != nil {
			log.Fatalf("Count of lines scan error: %s", err.Error())
		}

		checkDir, err := readJSONLines(in, lineCount)
		if err != nil {
			log.Fatalf("read JSON lines error: %s", err.Error())
		}

		count := checkAllDirsAndFiles(checkDir, false)

		fmt.Fprintln(out, count)
	}

}

func readJSONLines(in *bufio.Reader, linesCount int) (*CheckDir, error) {
	bytesJSON := make([]byte, 0)

	for i := 0; i <= linesCount; i++ {
		bytes, err := in.ReadBytes('\n')
		if err != nil {
			log.Fatalf("read bytes error: %s\n", err.Error())
		}

		bytesJSON = slices.Insert(bytesJSON, len(bytesJSON), bytes...)
	}

	directory := new(CheckDir)

	err := json.Unmarshal(bytesJSON, &directory)
	if err != nil {
		log.Fatalf("unmarshall error: %s\n", err.Error())
	}

	return directory, nil
}

func checkAllDirsAndFiles(checkDir *CheckDir, isParentInfected bool) int {
	hasInfectedFiles := isParentInfected

	infectedFilesCount := 0

	if hasInfectedFiles {
		//log.Printf("parent dir infected. all %d files in dir '%s' are infected", len(checkDir.Files), checkDir.Name)
		infectedFilesCount += len(checkDir.Files)
	} else {
		for i := 0; i < len(checkDir.Files); i++ {
			file := checkDir.Files[i]

			if strings.HasSuffix(file, ".hack") {
				hasInfectedFiles = true
				infectedFilesCount += len(checkDir.Files)
				//log.Printf("found infected file '%s'. all %d files in dir '%s' are infected", file, len(checkDir.Files), checkDir.Name)
				break
			}
		}
	}

	if len(checkDir.Folders) > 0 {
		for i := 0; i < len(checkDir.Folders); i++ {
			infectedFilesCount += checkAllDirsAndFiles(checkDir.Folders[i], hasInfectedFiles)
		}
	}

	return infectedFilesCount
}
