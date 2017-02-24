package main

//go:generate go-bindata -pkg assets -o generated/assets/assets.go -prefix assets/ assets/...
//go:generate swagger generate client --target generated --spec ./swagger.yml
import (
	"fmt"
	"github.com/bozaro/tech-db-forum/tests"
	"github.com/mkideal/cli"
	"github.com/op/go-logging"
	"net/url"
	"os"
	"reflect"
	"regexp"
)

type parserUrl struct {
	ptr interface{}
}

type parserRegexp struct {
	ptr interface{}
}

func newParserUrl(ptr interface{}) cli.FlagParser {
	return &parserUrl{ptr}
}

func newParserRegexp(ptr interface{}) cli.FlagParser {
	return &parserRegexp{ptr}
}

func (parser *parserUrl) Parse(s string) error {
	u, err := url.Parse(s)
	if err == nil {
		val := reflect.ValueOf(parser.ptr)
		val.Elem().Set(reflect.ValueOf(*u))
	}
	return err
}

func (parser *parserRegexp) Parse(s string) error {
	u, err := regexp.Compile(s)
	if err == nil {
		val := reflect.ValueOf(parser.ptr)
		val.Elem().Set(reflect.ValueOf(*u))
	}
	return err
}

type CmdCommonT struct {
	cli.Helper
	Url *url.URL `cli:"u,url" usage:"base url for testing API" parser:"url" dft:"http://localhost:5000/api"`
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
	Keep   bool           `cli:"k,keep" usage:"Don't stop after first failed test'"`
	Test   *regexp.Regexp `cli:"t,tests" usage:"Mask for running test names (regexp)" parser:"regexp" dft:".*"`
	Report string         `cli:"r,report" usage:"Detailed report file" dft:"report.html"`
}

var cmdFunc = &cli.Command{
	Name: "func",
	Desc: "run functional testing",
	Argv: func() interface{} { return new(CmdFuncT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CmdFuncT)
		os.Exit(tests.Run(argv.Url, argv.Test, argv.Report, argv.Keep))
		return nil
	},
}

func main() {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.4s}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	// Set the backends to be used.
	logging.SetBackend(logging.NewBackendFormatter(backend, format))

	cli.RegisterFlagParser("url", newParserUrl)
	cli.RegisterFlagParser("regexp", newParserRegexp)

	if err := cli.Root(root,
		cli.Tree(cmdFunc),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
