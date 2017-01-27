package main

//go:generate swagger generate client --target . --spec ./swagger.yml
import (
	"github.com/bozaro/tech-db-forum/tests"
)

func main() {
	tests.Run()
}
