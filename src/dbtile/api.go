package dbtile

import (
	"errors"
	"log"
	//"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MapMeta base Metadat info
type MapMeta struct {
	Maxzoom      int64  `json:"maxzoom"`
	Minzoom      int64  `json:"minzoom"`
	EtagPrefix   string `json:"etagprefix"`
	LastModified string `json:"lastmodified"`
}

// To store meta info used to serve all requests
var metaInfo MapMeta

// InitMeta init meta attribute :
// This is necessary to init values used to generate etags and last-modified header
func InitMeta() {
	// Get file metadata date
	_, metaDb, _ := DbGetMeta()
	metaInfo.EtagPrefix = metaDb.Mtime + "." + metaDb.Filesize + "." + metaDb.Format

	// Get meta min/max zoom
	tmpMaxzoom, tmpMaxzoomErr := strconv.ParseInt(metaDb.Maxzoom, 10, 64)
	if nil == tmpMaxzoomErr {
		metaInfo.Maxzoom = tmpMaxzoom
	} else {
		metaInfo.Maxzoom = 65000 // A big default Maxzoom to not have error "Out Of Bound"
	}
	tmpMinzoom, tmpMinzoomErr := strconv.ParseInt(metaDb.Minzoom, 10, 64)
	if nil == tmpMinzoomErr {
		metaInfo.Minzoom = tmpMinzoom
	} else {
		metaInfo.Minzoom = 0 // A smallest possible default Minzoom to not have error "Out Of Bound"
	}

	// Set last modified to meta info from files
	timsp, timspErr := strconv.ParseInt(metaDb.Mtime, 10, 64)
	if nil == timspErr {
		tm := time.Unix(timsp/1000, 0)
		metaInfo.LastModified = strings.Replace(tm.Format(time.RFC1123), "CEST", "GMT", 1) // a lasmodified header : https://developer.mozilla.org/fr/docs/Web/HTTP/Headers/Last-Modified
	} else {
		// No info from file : set LastModified to server start
		metaInfo.LastModified = strings.Replace(time.Now().Format(time.RFC1123), "CET", "GMT", 1)
	}
}

// -------------------------------------
// Public functions
// -------------------------------------

// WriteHTTPAnswer http answer to a tile request
func WriteHTTPAnswer(w http.ResponseWriter, r *http.Request, route string, XYZtoEPSG bool) {
	// Decode URL -> scheme is /ROUTE/z/x/y.pbf
	rStr := strings.TrimSuffix(strings.TrimPrefix(r.RequestURI, route), ".pbf")
	xyz := strings.Split(rStr, "/")
	// fmt.Printf("RequestURI: %+v\n", r.RequestURI) // DEBUG
	// fmt.Printf("z/x/y: %+v -> %d\n", xyz, len(xyz)) //DEBUG

	// Check there is 3 values
	if 3 != len(xyz) {
		http.Error(w, "Tile does not exist", 204)
	} else {
		// Get z/x/y values
		z, zErr := strconv.ParseInt(xyz[0], 10, 64)
		x, xErr := strconv.ParseInt(xyz[1], 10, 64)
		y, yErr := strconv.ParseInt(xyz[2], 10, 64)

		// Check z/x/y are integer
		if (nil != zErr) || (nil != xErr) || (nil != yErr) {
			http.Error(w, "Tile does not exist", 204) // Same message as tileserver-gl
		} else {
			// fmt.Printf("LEN=%s ID=%s \n", strconv.Itoa(len(m.Images_tile_data)), m.Map_tile_id)

			// Get tile
			tile, httpCode, tileErr := GetTile(x, y, z, XYZtoEPSG)
			if 200 != httpCode {
				if nil != tileErr {
					http.Error(w, tileErr.Error(), httpCode)
				} else {
					http.Error(w, "unknow error", httpCode)
				}
			}

			// return tile OK
			contentLength := strconv.Itoa(len(tile))
			w.Header().Set("Content-Length", contentLength)                                                          // set content length to avoid chuncked transfer
			w.Header().Set("ETag", `W/"`+xyz[0]+"."+xyz[1]+"."+xyz[2]+"-"+contentLength+"."+metaInfo.EtagPrefix+`"`) // Weak ETag
			w.Header().Set("Last-Modified", metaInfo.LastModified)                                                   // Last-modified is very important for browser to use cache ... and cache is important to limit a bit tile serving
			if strings.Contains(http.DetectContentType(tile), "gzip") {
				w.Header().Set("Content-Encoding", "gzip") // or we could check if byte start with "0x1f, 0x8b"
			}
			w.Header().Set("Content-Type", "application/x-protobuf") // We always serve protobuf content
			w.Write(tile)
		}
	}
}

// GetTile Get 1 tile according to x,y,z and XYZtoEPSG conversion
func GetTile(x int64, y int64, z int64, XYZtoEPSG bool) (tile []byte, httpCode int, err error) {
	// Check y, y, z
	if (z < metaInfo.Minzoom) || (z > metaInfo.Maxzoom) || (x < 0) || (y < 0) {
		log.Printf("Out of bound (%d/%d/%d) (XYZtoEPSG:%t)\n", z, x, y, XYZtoEPSG)

		httpCode = 204
		err = errors.New("Out of bounds") // Same message as tileserver-gl
	} else if XYZtoEPSG && ((float64(x) >= math.Pow(2, float64(z))) || (float64(y) >= math.Pow(2, float64(z)))) {
		log.Printf("Out of bound (%d/%d/%d) (XYZtoEPSG:%t)\n", z, x, y, XYZtoEPSG)

		httpCode = 204
		err = errors.New("Out of bounds") // Same message as tileserver-gl
	} else {
		// Fix y : needed to make it work with Spherical Mercator projection (projection code EPSG:3857) https://github.com/mapbox/sphericalmercator
		if XYZtoEPSG {
			y = int64(math.Pow(2, float64(z)) - 1 - float64(y))
		}

		// Get tile at x/y/z
		n := 0
		n, tile, err = DbGet1Tile(z, x, y)
		if (nil != err) || (n <= 0) {
			log.Printf("Tile (%d/%d/%d) does not exist -> %d / %s\n", x, y, z, n, err)

			httpCode = 204
			err = errors.New("Tile does not exists") // Same message as tileserver-gl
		} else {
			httpCode = 200 // Tile ok
		}
	}

	// Return tile
	return tile, httpCode, err
}
