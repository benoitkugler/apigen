package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/benoitkugler/apigen/fetch"
)

func main() {
	source := flag.String("source", "", "go source file containing the API")
	out := flag.String("out", "", "ts output file")
	typesFile := flag.String("types", "", "path to types declaration to import")
	flag.Parse()

	apis := fetch.FetchAPIs(*source)
	code := apis.Render()
	if *typesFile != "" {
		code = fmt.Sprintf("import * as types from %q\n", *typesFile) + code
	}
	if err := ioutil.WriteFile(*out, []byte(code), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
