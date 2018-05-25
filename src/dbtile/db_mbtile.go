package dbtile

import (
	//"fmt"
	"database/sql"
	"log"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DbMapMeta Metadata info
type DbMapMeta struct {
	Center    string `db:"center" json:"center"`
	Bounds    string `db:"bounds" json:"bounds"`
	Maxzoom   string `db:"maxzoom" json:"maxzoom"`
	Minzoom   string `db:"minzoom" json:"minzoom"`
	Mtime     string `db:"mtime" json:"mtime"`
	MaskLevel string `db:"maskLevel" json:"masklevel"`
	Format    string `db:"format" json:"format"`
	Filesize  string `db:"filesize" json:"filesize"`
	Type      string `db:"type" json:"type"`
}

// DbMap is map struct that match SQL select
type DbMap struct {
	ZoomLevel  int64  `db:"zoom_level" json:"zoom_level"`
	TileColumn int64  `db:"tile_column" json:"tile_column"`
	TileRow    int64  `db:"tile_row" json:"tile_row"`
	TileID     string `db:"tile_id" json:"tile_id"`
	TileData   []byte `db:"tile_data" json:"tile_data,omitempty"` //NOTE: the []byte will be base64 encoded by JSON Marshall (nice)
}

// mapSelectBase SQL select base (select maximum infos for 1 tile)
const mapSelectBase = `SELECT
m."zoom_level",
m."tile_column",
m."tile_row",
m."tile_id",
i."tile_data"
FROM "map" m
LEFT JOIN "images" i
	ON i."tile_id" = m."tile_id"` // Need to add 'limit 1'

// -------------------------------------
// Public functions
// -------------------------------------

// DbGet1Tile get tile data (only pbf image)
func DbGet1Tile(z int64, x int64, y int64) (int, []byte, error) {
	const query = `SELECT "tile_data" FROM "tiles" WHERE "zoom_level"=? AND "tile_column"=? AND "tile_row"=?`
	var tileData []byte
	nb := 0
	rows, rowsErr := dbConn.Query(query, z, x, y)
	if nil != rowsErr {
		return nb, tileData, rowsErr
	}
	defer rows.Close()

	// Scan all result and append it
	for rows.Next() {
		nb++
		var tmpTileData []byte
		rows.Scan(&tmpTileData) //tile_data blob
		tileData = append(tileData, tmpTileData...)
	}
	rows.Close()

	return nb, tileData, rowsErr
}

// DbGet1TileDetail get tile image + data by x/y coordinate and z zoom level
func DbGet1TileDetail(z int64, x int64, y int64) (int, *DbMap, error) {
	const query = mapSelectBase + ` WHERE m."zoom_level" = :z_level
	AND m."tile_column" = :x_coord
	AND m."tile_row" = :y_coord
	LIMIT 1`
	nb := 0
	dbMap := DbMap{}

	// Build query arguments
	queryArgs := map[string]interface{}{ // Here order don't matter
		"z_level": z,
		"x_coord": x,
		"y_coord": y,
	}

	// Prepare get query
	nstmt, err := dbConn.PrepareNamed(query)
	if nil != err {
		//fmt.Printf("ERR %s\n", err)
		return nb, &dbMap, err
	}
	// Do query
	err = nstmt.Get(&dbMap, queryArgs)

	if sql.ErrNoRows == err {
		err = nil // This is not an error, just a 0 return query
	} else if nil == err {
		nb = 1 // This query can only return 1 row
	} else {
		// fmt.Printf("ERR: %s\n", err)
	}

	return nb, &dbMap, err
}

// DbGetMeta get metadata
// Note: metadata are not mandatory and can not exist in an mtile sqlite file, so this function maxbe useless
func DbGetMeta() (int, *DbMapMeta, error) {
	const queryCenter = `SELECT "value" FROM "metadata" WHERE name="center" LIMIT 1`       // Ex: "8.5388,47.3434,128.5388,47.3434,12"
	const queryBounds = `SELECT "value" FROM "metadata" WHERE name="bounds" LIMIT 1`       // Ex: "-180,-85.05112877980659,180,85.0511287798066"
	const queryMaxzoom = `SELECT "value" FROM "metadata" WHERE name="maxzoom" LIMIT 1`     // Ex: "14"
	const queryMinzoom = `SELECT "value" FROM "metadata" WHERE name="minzoom" LIMIT 1`     // Ex: "0"
	const queryMtime = `SELECT "value" FROM "metadata" WHERE name="mtime" LIMIT 1`         // Ex: "1463000297761"
	const queryMaskLevel = `SELECT "value" FROM "metadata" WHERE name="maskLevel" LIMIT 1` // Ex: "8"
	const queryFormat = `SELECT "value" FROM "metadata" WHERE name="format" LIMIT 1`       // Ex: "pbf"
	const queryFilesize = `SELECT "value" FROM "metadata" WHERE name="filesize" LIMIT 1`   // Ex: "55719548"
	const queryType = `SELECT "value" FROM "metadata" WHERE name="type" LIMIT 1`           // Ex: "baselayer"

	meta := DbMapMeta{}
	nb := 0
	var err error
	err = dbConn.Get(&meta.Center, queryCenter)
	if nil != err {
		log.Printf("Meta err Center %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Bounds, queryBounds)
	if nil != err {
		log.Printf("Meta err Bounds %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Maxzoom, queryMaxzoom)
	if nil != err {
		log.Printf("Meta err Maxzoom %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Minzoom, queryMinzoom)
	if nil != err {
		log.Printf("Meta err Minzoom %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Mtime, queryMtime)
	if nil != err {
		log.Printf("Meta err Mtime %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.MaskLevel, queryMaskLevel)
	if nil != err {
		log.Printf("Meta err Masklevel %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Format, queryFormat)
	if nil != err {
		log.Printf("Meta err Format %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Filesize, queryFilesize)
	if nil != err {
		log.Printf("Meta err Filesize %s", err)
	} else {
		nb++
	}

	err = dbConn.Get(&meta.Type, queryType)
	if nil != err {
		log.Printf("Meta err Type %s", err)
	} else {
		nb++
	}

	return nb, &meta, err
}
