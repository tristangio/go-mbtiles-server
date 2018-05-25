package main

import (
	"dbtile"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Main server
// start a static file server to serve static file (demo_public/)
// instentiate db and wire it to dbtile
func main() {
	// File and line when a log occur
	log.SetFlags(log.Lshortfile)

	// Read arg
	// ----------------

	flagPort := flag.Int("port", 8086, "Port to use, effective only for non-https (default: 8086)")
	flagMbtiles := flag.String("mbtiles", "demo_mbtiles/2017-07-03_albania_tirana.mbtiles", "mbtiles file to use")
	flagAssetDir := flag.String("asset_dir", "demo_public/", "Set public file directory to serve (Browser files like mapboxgl.js, fonts, style etc, empty to not serve any browser files)")
	flagStyleInDir := flag.String("style_in", "demo_public/styles_src/bright/", "style to use : will replace {{{HOST}}} with your current IP")
	flagStyleOutDir := flag.String("style_out", "demo_public/styles/bright/", "where to store the replaced style_in")
	flagStyleHostName := flag.String("style_hostname", "", "Use a hostname for your style (will replace {{{HOST}}} with this value")
	flagRouteTile := flag.String("route_tile", "/map/", "Route browser will use to get tiles to display")
	flagXYZtoEPSG := flag.Bool("XYZ_to_EPSG", true, "Convert x/y/z to EPSG:3857 (true for Mapbox mbtiles usage)")
	flagServerReadTimeout := flag.Int("server_read_timeout", 40, "Server read request timeout in seconds")
	flagServerWriteTimeout := flag.Int("server_write_timeout", 40, "Server write request timeout in seconds")
	flagVerbose := flag.Bool("verbose", false, "true to enable verbose output")
	flag.Parse()
	portStr := fmt.Sprintf(":%d", *flagPort)

	// Database connection init
	// ----------------------

	dbErr := dbtile.InitDb(*flagMbtiles)
	if dbErr != nil {
		// Error in database init
		log.Fatal("Main : Database error : ", dbErr.Error())
	} else {
		// All is okay, better close the database before main exit (should happen only when dev)
		defer dbtile.CloseDb()
	}

	// Init meta info : Must be done before serving map
	dbtile.InitMeta()

	// Styles init
	// ----------------------

	if ("" != *flagAssetDir) && ("" != *flagStyleInDir) {
		if "" == *flagStyleHostName {
			InitStyle(*flagStyleInDir, *flagStyleOutDir, "http", GetOutboundIP(), portStr, *flagVerbose)
		} else {
			InitStyle(*flagStyleInDir, *flagStyleOutDir, "http", *flagStyleHostName, portStr, *flagVerbose)
		}
	}

	// Public file server (for login page)
	// ----------------------

	// Set a specific handler mux (it's better if you import other lib to not polute the default http handler, see: https://blog.cloudflare.com/exposing-go-on-the-internet/)
	h := http.NewServeMux()

	// Static file server (to serve login HTML + JS + CSS)
	// Here we use a file server that do not list files
	if "" != *flagAssetDir {
		fs := justFilesFilesystem{http.Dir(*flagAssetDir)} // Same as fs:=http.FileServer(http.Dir("../../webapp_public")) but not listing directory
		h.Handle("/", http.FileServer(fs))
	}

	// API route
	// ----------------------

	// API mbtile
	h.HandleFunc(*flagRouteTile, func(w http.ResponseWriter, r *http.Request) {
		dbtile.WriteHTTPAnswer(w, r, *flagRouteTile, *flagXYZtoEPSG)
	})

	// Start public static file server + API
	// ----------------------

	// Set a http server with timeout, do not use default httpserver (see: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
	srv := &http.Server{
		Addr:         portStr,
		ReadTimeout:  time.Duration(*flagServerReadTimeout) * time.Second,
		WriteTimeout: time.Duration(*flagServerWriteTimeout) * time.Second,
		Handler:      h,
	}

	// Starting plain text HTTP server
	fmt.Printf("Server Up and running on %s%s using plain/text http \n", GetOutboundIP(), portStr)
	srv.ListenAndServe()
}
