package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/benoitkugler/apigen/fetch"
	"github.com/benoitkugler/apigen/gents"
)

func main() {
	source := flag.String("source", "", "go source file containing the API")
	out := flag.String("out", "", "ts output file")
	flag.Parse()

	apis := fetch.FetchAPIs(*source)
	code := gents.Service(apis).Render()
	if err := ioutil.WriteFile(*out, []byte(code), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
