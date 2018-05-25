go-mtile-server
-----------------

A very simple server ready to serve your [mbtiles](https://wiki.openstreetmap.org/wiki/MBTiles) file.  
By default a mapbox-gl-js interface is used.

Note: This is more a demo on how to use the very simple mbtiles format in Go than a robust server.

If you want a more advanced server still in Go you should look at [consbio/mbtileserver](https://github.com/consbio/mbtileserver/blob/master/mbtiles/mbtiles.go)

And if you need to support old browser and generate raster tile from your vector tiles you should look at [tileserverGL](http://tileserver.org/) (NodeJS based)

How to build
-----------------

run :
```sh
$ ./set_go_path.sh
```
export your GOPATH as displayd, then :
```sh
$ go build
```

How to run
-----------------

You need an mbtiles file and the associated style.  
The easy way is to download vectortile from openmaptile : download [andorra for exemple](https://openmaptiles.com/downloads/dataset/osm/europe/andorra/#10.14/42.5425/1.5999).  
You can also generate your own vector tile with [openmaptile](https://github.com/openmaptiles/openmaptiles) or [tippecanoe](https://github.com/mapbox/tippecanoe)

Put your mbtiles file somewhere and run :
```sh
$ ./gomtileserver -mbtiles path/to/your/file.mbtiles
```

By default style [osm-bright](https://github.com/openmaptiles/osm-bright-gl-style) is used. To change it you need to change the path "style" in demo_public/index.html

For ease gomtileserver transform style in demo_public/style_src/bright and generate one with current server IP (by changing the {{{HOST}}} value).
TODO: do it better or allow to overwrite the HOST
