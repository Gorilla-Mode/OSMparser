# OSMparser
Parses OSM JSON to postGIS point geometries

## Usage
1. Drop JSON files into `./in`
2. Compile and run `main.go`'
3. Parsed SQL files will be in `./out`

##  Exceptions
###  Input
- OSM JSON file

The name of the file will be the "type" value in the parsed SQL file.

### Output format

| **Output**       | **id**                | **type**                                            | **geom**      |
|------------------|-----------------------|-----------------------------------------------------|---------------|
| **Description**  | ID from OSM JSON file | From filename. e.g `chemist.json` returns `chemist` | PostGIS point |
| **Default name** | fid                   | key                                                 | wkt_geom      |

#### Ouput example
```sql
INSERT INTO t (fid, key, wkt_geom) VALUES (1::bigint, 'myballs', point(1.0, 2.0)) ON CONFLICT DO NOTHING;
```

## Geometries
 ### Support
 - [x] Nodes
 - [x] Ways
 - [ ] Multipolygons

> [!IMPORTANT]  
> Geometry returned for a way is the centroid of the way.