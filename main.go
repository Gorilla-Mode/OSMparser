package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	s "strings"
)

type Element struct {
	Type string            `json:"type"`
	ID   int64             `json:"id"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags"`
}

type OSMData struct {
	Elements []Element `json:"elements"`
}

const (
	Node  string = "node"
	Way          = "way"
	Multi        = "multipolygon"
)

func main() {
	inDir := "in"
	outDir := "out"

	files, err := os.ReadDir(inDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}

		processFile(inDir, outDir, file.Name())
	}
}

func processFile(inDir, outDir, fileName string) {
	data, err := os.ReadFile(inDir + "/" + fileName)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileName, err)
		return
	}

	// Remove BOM if present
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var osm OSMData
	if err := json.Unmarshal(data, &osm); err != nil {
		fmt.Printf("Error unmarshaling %s: %v\n", fileName, err)
		return
	}

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return
	}

	outFileName := s.TrimSuffix(fileName, ".json") + ".sql"
	out, err := os.Create(filepath.Join(outDir, outFileName))
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", outFileName, err)
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	_, err = out.WriteString("BEGIN;\n")
	if err != nil {
		return
	}
	count := 0
	buildingType := s.TrimSuffix(fileName, ".json")

	for _, elem := range osm.Elements {
		if elem.Tags == nil {
			continue
		}

		switch elem.Type {
		case Node:
			_, done := parseNode(elem, err, out, buildingType, &count)
			if done {
				break
			}
		case Way:
			fmt.Printf("Feat not implemented yet: %s\n", elem.Type)
			break

		case Multi:
			fmt.Printf("Feat not implemented yet: %s\n", elem.Type)
			break
		}
	}

	_, err = out.WriteString("COMMIT;\n")
	if err != nil {
		return
	}
	fmt.Printf("Parsed file %s, found %d building nodes\n", buildingType, count)
}

func parseNode(elem Element, err error, out *os.File, buildingType string, count *int) (error, bool) {

	_, err = out.WriteString(fmt.Sprintf(
		"INSERT INTO buildings (fid, key, wkt_geom) VALUES (%d::bigint, '%s', point(%f, %f));\n",
		elem.ID, buildingType, elem.Lon, elem.Lat,
	))
	if err != nil {
		return nil, true
	}
	*count++
	return err, false
}
