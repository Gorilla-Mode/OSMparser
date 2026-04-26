package main

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
