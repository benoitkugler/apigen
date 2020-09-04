package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/benoitkugler/apigen/fetch"
)

func main() {
	source := flag.String("source", "", "go source file containing the API")
	out := flag.String("out", "", "ts output file")
	flag.Parse()

	apis := fetch.FetchAPIs(*source)
	code := apis.Render()
	if err := ioutil.WriteFile(*out, []byte(code), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
