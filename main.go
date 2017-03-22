package main

//go:generate go-bindata -nometadata -pkg assets -o generated/assets/assets.go -prefix assets/ assets/...
//go:generate swagger generate client --target generated --spec ./swagger.yml
import (
	"fmt"
	"github.com/bozaro/tech-db-forum/tests"
	"github.com/mkideal/cli"
	"github.com/op/go-logging"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"time"
)

const (
	EXIT_INVALID_COMMAND = iota + 1
	EXIT_WAIT_ALIVE_TIMEOUT
	EXIT_FUNC_FAILED
	EXIT_FILL_FAILED
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
	Url       *url.URL `cli:"u,url" usage:"base url for testing API" parser:"url" dft:"http://localhost:5000/api"`
	WaitAlive int      `cli:"wait" usage:"wait before remote API make alive (while connection refused or 5XX error on base url)" dft:"30"`
}

var root = &cli.Command{
	Desc: "https://github.com/bozaro/tech-db-forum",
	Argv: func() interface{} { return nil },
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		os.Exit(EXIT_INVALID_COMMAND)
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
		waitAlive(argv.CmdCommonT)
		if tests.Run(argv.Url, argv.Test, argv.Report, argv.Keep) > 0 {
			os.Exit(EXIT_FUNC_FAILED)
		}
		return nil
	},
}

type CmdFillT struct {
	CmdCommonT
}

var cmdFill = &cli.Command{
	Name: "fill",
	Desc: "fill database with random data",
	Argv: func() interface{} { return new(CmdFillT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CmdFillT)
		waitAlive(argv.CmdCommonT)
		if tests.Fill(argv.Url) > 0 {
			os.Exit(EXIT_FILL_FAILED)
		}
		return nil
	},
}

var cmdVersion = &cli.Command{
	Name: "version",
	Desc: "show version",
	Argv: func() interface{} { return new(CmdFillT) },
	Fn: func(ctx *cli.Context) error {
		fmt.Println(tests.VersionFull())
		return nil
	},
}
var log = logging.MustGetLogger("main")

func waitAlive(argv CmdCommonT) {
	req, err := http.NewRequest("GET", argv.Url.String(), nil)
	if err != nil {
		panic(err)
	}
	lst := ""

	if argv.WaitAlive <= 0 {
		return
	}
	timeout := time.Now().Add(time.Duration(argv.WaitAlive) * time.Second)
	for time.Now().Before(timeout) {
		msg := ""
		if err == nil {
			res, err := tests.HttpTransport.RoundTrip(req)
			if err != nil {
				msg = fmt.Sprintf("Connection error: %s", err.Error())
			} else if res.StatusCode >= 500 && res.StatusCode < 600 {
				msg = fmt.Sprintf("Invalid response code: %d", res.StatusCode)
			} else {
				if lst != "" {
					log.Info("Service is alive")
				}
				return
			}
		}
		if lst != msg {
			log.Warning("Service unavailable: " + msg)
			lst = msg
		}
		time.Sleep(time.Second / 10)
	}
	log.Error("Wait service alive timeout")
	os.Exit(EXIT_WAIT_ALIVE_TIMEOUT)
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
		cli.Tree(cmdFill),
		cli.Tree(cmdVersion),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_INVALID_COMMAND)
	}
}
