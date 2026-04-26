package main

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
