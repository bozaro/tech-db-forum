package main

//go:generate swagger generate client --target generated --spec ./swagger.yml
import (
	"fmt"
	"github.com/bozaro/tech-db-forum/tests"
	"github.com/mkideal/cli"
	"net/url"
	"os"
	"reflect"
)

type parserUrl struct {
	ptr interface{}
}

func newParserUrl(ptr interface{}) cli.FlagParser {
	return &parserUrl{ptr}
}

func (parser *parserUrl) Parse(s string) error {
	u, err := url.Parse(s)
	if err == nil {
		val := reflect.ValueOf(parser.ptr)
		val.Elem().Set(reflect.ValueOf(*u))
	}
	return err
}

type CmdCommonT struct {
	cli.Helper
	Url *url.URL `cli:"url" usage:"base url for testing API" parser:"url" dft:"http://localhost:5000/api"`
}

var root = &cli.Command{
	Desc: "https://github.com/bozaro/tech-db-forum",
	Argv: func() interface{} { return nil },
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		os.Exit(1)
		return nil
	},
}

type CmdFuncT struct {
	CmdCommonT
}

var cmdFunc = &cli.Command{
	Name: "func",
	Desc: "run functional testing",
	Argv: func() interface{} { return new(CmdFuncT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CmdFuncT)
		os.Exit(tests.Run(argv.Url))
		return nil
	},
}

func main() {
	cli.RegisterFlagParser("url", newParserUrl)

	if err := cli.Root(root,
		cli.Tree(cmdFunc),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
