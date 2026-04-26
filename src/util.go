package main

import (
	"fmt"
	"os"
	s "strings"
)

const (
	Reset = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	Gray
	White
)

var ansi = map[int]string{
	Reset:   "\033[0m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m,",
	Cyan:    "\033[36m",
	Gray:    "\033[37m",
	White:   "\033[97m",
}

func calculateCentroid(geometry []Coordinate) (lat, lon float64) {
	if len(geometry) == 0 {
		return 0, 0
	}

	for _, coord := range geometry {
		lat += coord.Lat
		lon += coord.Lon
	}

	lat /= float64(len(geometry))
	lon /= float64(len(geometry))

	return lat, lon
}

func printResults(results []ParseResult) {
	for i, result := range results {
		prefix := "├─"
		connector := "│\t"
		if i == len(results)-1 {
			prefix = "└─"
			connector = "\t"
		}
		fmt.Printf("\t%s Parsed file %s in %s\n\t%s├── nodes: %s%d%s\n\t%s├── way(s): %s%d%s\n\t%s└── multipolygon(s): %s%d%s\n",
			prefix, result.fileName, result.elapsed, connector, ansi[Green], result.nodes, ansi[Reset],
			connector, ansi[Green], result.ways, ansi[Reset], connector, ansi[Green], result.multipolygons, ansi[Reset])
	}
}

func printFiles(countFiles int, inPath string, files []os.DirEntry) {
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
}
