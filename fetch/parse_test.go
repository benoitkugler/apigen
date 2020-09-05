package fetch

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	pack, file, err := loadSource("test/routes.go")
	// pack, file, err := loadSource("/home/benoit/go/src/github.com/benoitkugler/goACVE/server/main.go")
	// pack, file, err := loadSource("/home/benoit/go/src/github.com/benoitkugler/intendance/server/main.go")
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
