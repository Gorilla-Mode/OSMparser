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

func parseRelation(elem Element, err error, out *os.File, buildingType string, count *int) (error, bool) {
	if elem.Bounds == nil {
		return nil, true
	}

	geometry := []Coordinate{
		{Lat: elem.Bounds.MinLat, Lon: elem.Bounds.MinLon},
		{Lat: elem.Bounds.MaxLat, Lon: elem.Bounds.MaxLon},
	}

	centerLat, centerLon := calculateCentroid(geometry)

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

func parseFile(inDir, outDir, fileName string, isLast bool) {
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

	delegateParser(osm, err, out, buildingType, &countNode, &countWay, &countMulti)

	_, err = out.WriteString("COMMIT;\n")
	if err != nil {
		return
	}

	end := t.Now()
	elapsed := end.Sub(start)
	if !isLast {
		fmt.Printf("\t├─┬─ Parsed file %s in %s\n\t│ └─┬── nodes: %s%d%s\n\t│   ├── way(s): %s%d%s\n\t"+
			"│   └── multipolygon(s): %s%d%s\n",
			fileName, elapsed, ansi[Green], countNode, ansi[Reset], ansi[Green], countWay, ansi[Reset], ansi[Green],
			countMulti, ansi[Reset])
	} else {
		fmt.Printf("\t└─┬─ Parsed file %s in %s\n\t  └─┬── nodes: %s%d%s\n\t    ├── way(s): %s%d%s\n\t"+
			"    └── multipolygon(s): %s%d%s\n\n",
			fileName, elapsed, ansi[Green], countNode, ansi[Reset], ansi[Green], countWay, ansi[Reset], ansi[Green],
			countMulti, ansi[Reset])
	}
}

func delegateParser(osm OSMData, err error, out *os.File, buildingType string, countNode *int, countWay *int, countMulti *int) {
	for _, elem := range osm.Elements {
		if elem.Tags == nil {
			continue
		}

		switch elem.Type {
		case Node:
			_, done := parseNode(elem, err, out, buildingType, countNode)
			if done {
				break
			}
		case Way:
			_, done := parseWay(elem, err, out, buildingType, countWay)
			if done {
				break
			}
		case Multi:
			_, done := parseRelation(elem, err, out, buildingType, countMulti)
			if done {
				break
			}
		}
	}
}
