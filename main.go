package main

//go:generate swagger generate client --target generated --spec ./swagger.yml
import (
	"github.com/bozaro/tech-db-forum/tests"
	"github.com/mkideal/cli"
	"net/url"
	"os"
)

type argT struct {
	cli.Helper
	ApiUrl string `cli:"url" usage:"base url for testing API" dft:"http://localhost:5000/api"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		u, err := url.Parse(argv.ApiUrl)
		if err != nil {
			return err
		}
		os.Exit(tests.Run(u))
		return nil
	})
}
