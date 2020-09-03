package fetch

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	pack, file, err := loadSource("test/routes.go")
	if err != nil {
		t.Fatal(err)
	}
	ti := time.Now()
	apis := parse(pack, file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(time.Since(ti))
	fmt.Println(apis)
}