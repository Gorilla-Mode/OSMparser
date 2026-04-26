package main

import (
	"fmt"
	"os"
	"path/filepath"
	s "strings"
	"sync"
	t "time"
)

var Cols = newColsNames("fid", "key", "wkt_geom", "buildings")

func main() {
	inDir := "in"
	outDir := "out"
	inPath, _ := filepath.Abs(inDir)
	files, err := os.ReadDir(inDir)
	if err != nil {
		panic(err)
	}
	countFiles := len(files)

	if countFiles == 0 {
		fmt.Println("No files found in: " + inPath)
		return
	}

	fmt.Printf("Found %d files in %s\n", countFiles, inPath)
	fmt.Printf("Staged files:\n")
	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}
		if files[countFiles-1] != file {
			fmt.Printf("\t├─ Found file: %s\n", file.Name())
		} else {
			fmt.Printf("\t└─ Found file: %s\n", file.Name())
		}

	}

	fmt.Printf("Parse? (y):")
	var i string
	_, err = fmt.Scan(&i)
	if err != nil {
		return
	}

	if i != "y" {
		return
	}

	fmt.Printf("Default cols? (y):")
	_, err = fmt.Scan(&i)
	if err != nil {
		return
	}

	if i != "y" {
		setColsNames()
	}

	fmt.Printf("\nParsing %d files...\n", countFiles)

	results := make([]ParseResult, 0, countFiles)
	var mu sync.Mutex
	var wg sync.WaitGroup

	start := t.Now()
	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}

		wg.Add(1)
		go func(f os.DirEntry) {
			defer wg.Done()
			result := parseFile(inDir, outDir, f.Name())
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(file)
	}
	wg.Wait()
	end := t.Now()
	elapsed := end.Sub(start)

	printResults(results)

	fmt.Printf("%sParsed %d files in %s%s\n", ansi[Green], countFiles, elapsed, ansi[Reset])
}
