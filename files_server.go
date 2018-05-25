package main

// Implement structure/filesystem that do not list directories content
// stolen here : https://groups.google.com/forum/#!msg/golang-nuts/bStLPdIVM6w/AXLz0hNqCrUJ
// Code from Brad FitzPatrick (one of the go lead dev, so it should be okay)

import (
	"net/http"
	"os"
	//"log"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	// Optionnal serve GZIP ... but it lacks setting header to tell it's gziped + checking client header for gzip support + it should check if zip is younger than original
	//var f http.File
	//var err error
	//if f, err = fs.fs.Open(name+".gz"); err == nil {
	//	log.Printf("Seving gziped %s \n", name)// There is a gziped version
	//} else {
	//	log.Printf("NO gziped %s \n", name_gz)
	//	f, err = fs.fs.Open(name) // no zip
	//}
	f, err := fs.fs.Open(name)
   if err != nil {
      return nil, err
   }
   return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// Use it like that :
//--------------------------

// func main() {
// 	fs := justFilesFilesystem{http.Dir("/tmp/")}
// 	http.ListenAndServe(":8080", http.FileServer(fs))
// }
