package main

//go:generate swagger generate client --target generated --spec ./swagger.yml
import (
	"github.com/bozaro/tech-db-forum/tests"
	"os"
)

func main() {
	os.Exit(tests.Run())
}
