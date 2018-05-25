package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// InitStyle list styles in input directory replace {{{HOST}}} by http://HOST:PORT and write new style in output directory
// Exemple:
// dirIn: demo_public/styles_src/bright
// dirOut: demo_public/styles/bright
// protocol: http
// host: pass your IP by using NetGetOutboundIP() from network_utils.go
// port: 8086
// verbose: true
func InitStyle(dirIn string, dirOut string, protocol string, host string, port string, verbose bool) {
	if verbose {
		fmt.Printf("Map style -> InitStyles from  %s --> %s \n", dirIn, dirOut)
	}

	files, filesErr := ioutil.ReadDir(dirIn)
	if nil != filesErr {
		log.Printf("File err %s\n", filesErr)
		return
	}
	for _, f := range files {

		// Read file
		input, err := ioutil.ReadFile(dirIn + f.Name())
		if err != nil {
			log.Printf("InitStyles -> err on file %s : %s\n", dirIn+f.Name(), err)
			continue
		}

		// Replace {{{HOST}}}
		output := string(input)
		if strings.HasSuffix(f.Name(), ".json") { // replace only json files (we don't want to mess sprite files or other)
			if !strings.HasSuffix(protocol, "://") {
				protocol += "://"
			}
			toReplace := protocol + host
			if !strings.HasPrefix(port, ":") {
				toReplace += ":"
			}
			toReplace += port
			output = strings.Replace(string(input), "{{{HOST}}}", toReplace, -1)
		}

		// Write file to new destination
		err = ioutil.WriteFile(dirOut+f.Name(), []byte(output), 0644)
		if err != nil {
			log.Printf("InitStyles -> err writing file %s : %s\n", dirOut+f.Name(), err)
			continue
		} else if verbose {
			log.Printf("Map style ->  %s : %s  -->  %s , OK\n", f.Name(), dirIn, dirOut)
		}
	}
}
