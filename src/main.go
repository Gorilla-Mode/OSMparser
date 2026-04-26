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

	printFiles(countFiles, inPath, files)

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

	results := make(chan ParseResult, countFiles)
	var wg sync.WaitGroup
	var elapsed t.Duration

	start := t.Now()
	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}

		wg.Add(1)
		go func(f os.DirEntry) {
			defer wg.Done()
			result := parseFile(inDir, outDir, f.Name())
			results <- result
		}(file)
	}
	go func() {
		wg.Wait()
		elapsed = t.Now().Sub(start)
		close(results)
	}()

	parsedResults := make([]ParseResult, 0, countFiles)
	for result := range results {
		parsedResults = append(parsedResults, result)
	}

	printResults(parsedResults)
	fmt.Printf("%sParsed %d files in %s%s\n", ansi[Green], countFiles, elapsed, ansi[Reset])
}
