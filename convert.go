package main

import (
	"log"
	"path/filepath"
	"flag"
	"strings"
	"fmt"
	"translation_converter/format"
	"encoding/json"
	"io/ioutil"
)

type Format interface {
	ToJson(string) error
	FromJson(string, map[string]string) error
	Ext() string
}

var formats = map[string]func(string) Format {
	"simple" : func(from string) Format {return format.Simple{Src:from}}}

func listFormats() {
	fmt.Println("Supported formats: ")
	for f, _ := range formats {
		fmt.Println("  * " + f)
	}
}

func main() {
	from := flag.String("from", "", "the path to the file to convert (REQUIRED)")
	to := flag.String("to", "", "the output path")
	help := flag.Bool("h", false, "Print usage menu")
	list := flag.Bool("l", false, "List formats")
	formatFlag := flag.String("format", "simple", "The format to convert from/to")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *list {
		listFormats()
		return
	}

	if *from == "" {
		log.Fatalf("The 'from' parameter is required")
	}

	formatObj := formats[*formatFlag](*from)
	fromJson := filepath.Ext(*from) == ".json"

	if *to == "" {
		var ext = "json"
		if (fromJson) {
			ext = formatObj.Ext()
		}
		*to = strings.Split(filepath.Base(*from), ".")[0]+"."+ext
	}

	fmt.Printf("Reading from %q\n", *from)
	fmt.Printf("Writing to %q\n", *to)

	var err error
	if fromJson {
		fmt.Printf("Converting from %q format to %q format\n", "json", *formatFlag)

		var jsonMap map[string]string
		jsonData, err := ioutil.ReadFile(*from)
		if err != nil {
			log.Fatalf("Error reading the json file %q due to: %v\n", *from, err)
		}
		json.Unmarshal(jsonData, &jsonMap)
		err = formatObj.FromJson(*to, jsonMap)
	} else {
		fmt.Printf("Converting from %q format to %q format\n", *formatFlag, "json")
		err = formatObj.ToJson(*to)
	}

	if err != nil {
		log.Fatalf("Failed conversion due to: %v\n", err)
	}


}

