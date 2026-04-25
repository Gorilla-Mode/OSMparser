package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	s "strings"
	t "time"
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Bounds struct {
	MinLat float64 `json:"minlat"`
	MinLon float64 `json:"minlon"`
	MaxLat float64 `json:"maxlat"`
	MaxLon float64 `json:"maxlon"`
}

type Element struct {
	Type string            `json:"type"`
	ID   int64             `json:"id"`
	Lat  float64           `json:"lat"`
	Lon  float64           `json:"lon"`
	Tags map[string]string `json:"tags"`
	// Way-specific fields
	Bounds   *Bounds      `json:"bounds,omitempty"`
	Nodes    []int64      `json:"nodes,omitempty"`
	Geometry []Coordinate `json:"geometry,omitempty"`
}

type OSMData struct {
	Elements []Element `json:"elements"`
}

const (
	Node  string = "node"
	Way          = "way"
	Multi        = "multipolygon"
)

type (
	tableDefinition struct {
		id    string
		btype string
		geom  string
		table string
	}
)

func newColsNames(id string, btype string, geom string, table string) *tableDefinition {
	return &tableDefinition{
		id:    id,
		btype: btype,
		geom:  geom,
		table: table,
	}
}

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
		fmt.Println("No files found in " + inPath)
		return
	}

	fmt.Printf("Found %d files in %s\n", countFiles, inPath)
	fmt.Printf("Staged files:\n")
	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}

		fmt.Printf("\tFound file: %s\n", file.Name())
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

	start := t.Now()
	for _, file := range files {
		if file.IsDir() || !s.HasSuffix(file.Name(), ".json") {
			continue
		}

		processFiles(inDir, outDir, file.Name())
	}
	end := t.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Parsed %d files in %s\n", countFiles, elapsed)
}

func processFiles(inDir, outDir, fileName string) {
	start := t.Now()
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
	buildingType := s.TrimSuffix(fileName, ".json")
	countNode := 0
	countWay := 0
	countMulti := 0

	for _, elem := range osm.Elements {
		if elem.Tags == nil {
			continue
		}

		switch elem.Type {
		case Node:
			_, done := parseNode(elem, err, out, buildingType, &countNode)
			if done {
				break
			}
		case Way:
			_, done := parseWay(elem, err, out, buildingType, &countWay)
			if done {
				break
			}
		case Multi:
			//TODO: Implement func
			break
		}
	}

	_, err = out.WriteString("COMMIT;\n")
	if err != nil {
		return
	}

	end := t.Now()
	elapsed := end.Sub(start)
	fmt.Printf("\tParsed file %s in %s\n\t\t-> nodes: %d\n\t\t-> way(s): %d\n\t\t-> multipolygon(s): %d\n",
		fileName, elapsed, countNode, countWay, countMulti)
}

func parseNode(elem Element, err error, out *os.File, buildingType string, count *int) (error, bool) {

	_, err = out.WriteString(fmt.Sprintf(
		"INSERT INTO buildings (%s, %s, %s) VALUES (%d::bigint, '%s', point(%f, %f)) ON CONFLICT DO NOTHING;\n",
		Cols.id, Cols.btype, Cols.geom, elem.ID, buildingType, elem.Lon, elem.Lat,
	))
	if err != nil {
		return nil, true
	}
	*count++
	return err, false
}

func parseWay(elem Element, err error, out *os.File, buildingType string, count *int) (error, bool) {
	if len(elem.Geometry) == 0 {
		return nil, true
	}

	centerLat, centerLon := calculateCentroid(elem.Geometry)

	_, err = out.WriteString(fmt.Sprintf(
		"INSERT INTO buildings (%s, %s, %s) VALUES (%d::bigint, '%s', point(%f, %f)) ON CONFLICT DO NOTHING;\n",
		Cols.id, Cols.btype, Cols.geom, elem.ID, buildingType, centerLon, centerLat,
	))
	if err != nil {
		return nil, true
	}
	*count++
	return err, false
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

func setColsNames() {
	fmt.Println("Set table name: ")
	var table string
	_, err := fmt.Scan(&table)
	if err != nil {
		return
	}

	fmt.Println("Set identity: ")
	var id string
	_, err = fmt.Scan(&id)
	if err != nil {
		return
	}

	fmt.Println("Set building type: ")
	var btype string
	_, err = fmt.Scan(&btype)
	if err != nil {
		return
	}

	fmt.Println("Set geometry column name: ")
	var geom string
	_, err = fmt.Scan(&geom)
	if err != nil {
		return
	}

	Cols.table = table
	Cols.id = id
	Cols.btype = btype
	Cols.geom = geom
}
