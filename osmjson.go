package main

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
	Multi        = "relation"
)
