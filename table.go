package main

import "fmt"

type tableDefinition struct {
	id    string
	btype string
	geom  string
	table string
}

func newColsNames(id string, btype string, geom string, table string) *tableDefinition {
	return &tableDefinition{
		id:    id,
		btype: btype,
		geom:  geom,
		table: table,
	}
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
