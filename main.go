package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/benoitkugler/apigen/fetch"
	"github.com/benoitkugler/structgen/enums"
)

func main() {
	source := flag.String("source", "", "go source file containing the API")
	out := flag.String("out", "", "ts output file")
	flag.Parse()

	pkg, f, err := fetch.LoadSource(*source)
	if err != nil {
		log.Fatalf("can't type check package : %s", err)
	}
	apis := fetch.Parse(pkg, f)

	enumTable, err := enums.FetchEnums(pkg)
	if err != nil {
		log.Fatalf("cant't parse enums : %s", err)
	}

	code := apis.Render(enumTable)

	if err := ioutil.WriteFile(*out, []byte(code), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	log.Printf("Api generated in %s", *out)
}
