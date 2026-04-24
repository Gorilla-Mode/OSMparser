# OSMparser
Parses OSM JSON to postGIS point geometries

## 1. Usage
1. Drop JSON files into `./in`
2. Compile and run `main.go`'
3. Parsed SQL files will be in `./out`

## 2. Exceptions
### 1. Table format
Parser expects a table with the following columns:

| fid     | type | wkt_geom |
|---------|------|----------|
| bingint | text | point    |

### 3. Input
- OSM JSON file

The name of the file will be the "type" value in the parsed SQL file.

### 4, Output
- PostGIS 