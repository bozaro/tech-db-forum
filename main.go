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
	Url             *url.URL `cli:"u,url" usage:"base url for testing API" parser:"url" dft:"http://localhost:5000/api"`
	WaitAlive       int      `cli:"wait" usage:"wait before remote API make alive (while connection refused or 5XX error on base url)" dft:"30"`
	DontCheckUpdate bool     `cli:"no-check" usage:"do not check version update"`
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
		commonPrepare(argv.CmdCommonT)
		if tests.Run(argv.Url, argv.Test, argv.Report, argv.Keep) > 0 {
			os.Exit(EXIT_FUNC_FAILED)
		}
		return nil
	},
}

type CmdFillT struct {
	CmdCommonT
	Threads   int    `cli:"t,thread" usage:"number of threads for generating data" dft:"8"`
	StateFile string `cli:"o,state" usage:"State file with information about database objects" dft:"tech-db-forum.dat"`
}

var cmdFill = &cli.Command{
	Name: "fill",
	Desc: "fill database with random data",
	Argv: func() interface{} { return new(CmdFillT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CmdFillT)
		commonPrepare(argv.CmdCommonT)
		perf := tests.Fill(argv.Url, argv.Threads, tests.NewPerfConfig())
		if perf == nil {
			os.Exit(EXIT_FILL_FAILED)
		}
		file, err := os.Create(argv.StateFile)
		if err != nil {
			log.Error("Can't create file: " + argv.StateFile)
			os.Exit(EXIT_FILL_FAILED)
		}
		err = perf.Save(file)
		if err != nil {
			log.Error("Can't save to file file: " + argv.StateFile)
			os.Exit(EXIT_FILL_FAILED)
		}
		return nil
	},
}

type CmdPerfT struct {
	CmdCommonT
	Threads int `cli:"t,thread" usage:"number of threads for performance testing" dft:"8"`
}

var cmdPerf = &cli.Command{
	Name: "perf",
	Desc: "run performance testing",
	Argv: func() interface{} { return new(CmdPerfT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CmdPerfT)
		commonPrepare(argv.CmdCommonT)
		perf := tests.Fill(argv.Url, argv.Threads, tests.NewPerfConfig())
		if perf == nil {
			os.Exit(EXIT_FILL_FAILED)
		}

		perf.Run(argv.Threads)
		return nil
	},
}

var cmdVersion = &cli.Command{
	Name: "version",
	Desc: "show version",
	Argv: func() interface{} { return new(CmdFillT) },
	Fn: func(ctx *cli.Context) error {
		fmt.Println(tests.VersionFull())
		if ver, err := tests.VersionCheck(); err == nil {
			switch ver {
			case tests.VERSION_LATEST:
				log.Infof("You use latest version of %s tool.", tests.Project)
			case tests.VERSION_LOCAL:
				log.Infof("You use local build of %s tool.", tests.Project)
			case tests.VERSION_OUTDATE:
				log.Warningf("You use outdated version of %s tool. Please update.", tests.Project)
			}
		}
		return nil
	},
}
var log = logging.MustGetLogger("main")

func commonPrepare(argv CmdCommonT) {
	checkUpdate(argv)
	waitAlive(argv)
}

func checkUpdate(argv CmdCommonT) {
	if !argv.DontCheckUpdate {
		if ver, _ := tests.VersionCheck(); ver == tests.VERSION_OUTDATE {
			log.Warningf("You use outdated version of %s tool. Please update.", tests.Project)
		}
	}
}

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
		cli.Tree(cmdPerf),
		cli.Tree(cmdVersion),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_INVALID_COMMAND)
	}
}
