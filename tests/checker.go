package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/assets"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/go-openapi/runtime"
	http_transport "github.com/go-openapi/runtime/client"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/op/go-logging"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
)

var log = logging.MustGetLogger("checker")

type Checker struct {
	// Имя текущей проверки.
	Name string
	// Описание текущей проверки.
	Description string
	// Функция для текущей проверки.
	FnCheck func(c *client.Forum, f *Factory)
	// Тесты, без которых проверка не имеет смысл.
	Deps []string
}

var s_templateUid int32 = 0

type CheckerByName []Checker

func (a CheckerByName) Len() int           { return len(a) }
func (a CheckerByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CheckerByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type CheckerTransport struct {
	t      runtime.ClientTransport
	report *Report
}

type CheckerClientResponseReader struct {
	reader runtime.ClientResponseReader
}

func EasyJSONConsumer() runtime.Consumer {
	return runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		if v, ok := data.(easyjson.Unmarshaler); ok {
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			l := jlexer.Lexer{Data: data}
			v.UnmarshalEasyJSON(&l)
			return l.Error()
		}
		dec := json.NewDecoder(reader)
		dec.UseNumber() // preserve number formats
		return dec.Decode(data)
	})
}

func (self CheckerClientResponseReader) Consume(r io.Reader, t interface{}) error {
	b := make([]byte, 1024)
	for {
		size, err := r.Read(b)
		if err != nil {
			return err
		}
		if size <= 0 {
			break
		}
	}
	return nil
}

func (self CheckerClientResponseReader) ReadResponse(response runtime.ClientResponse, _ runtime.Consumer) (interface{}, error) {
	return self.reader.ReadResponse(response, self)
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := NewValidator(operation.Context, self.report)
	operation.Client = &http.Client{Transport: tracker}
	if operation.Context != nil && operation.Context.Value(KEY_SKIP) != nil {
		operation.Reader = CheckerClientResponseReader{operation.Reader}
	}
	return self.t.Submit(operation)
}

func Checkpoint(c *client.Forum, message string) bool {
	return c.Transport.(*CheckerTransport).report.Checkpoint(message)
}

var registeredChecks []Checker

func Register(checker Checker) {
	registeredChecks = append(registeredChecks, checker)
}

func RunCheck(check Checker, report *Report, url *url.URL) {
	report.Result = Success
	transport := CreateTransport(url)
	defer func() {
		if r := recover(); r != nil {
			report.AddError(r)
		}
	}()
	check.FnCheck(client.New(&CheckerTransport{transport, report}, nil), NewFactory())
}

func CreateTransport(url *url.URL) *http_transport.Runtime {
	cfg := client.DefaultTransportConfig()
	if url != nil {
		cfg.WithHost(url.Host).WithSchemes([]string{url.Scheme}).WithBasePath(url.Path)
	}
	transport := http_transport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	transport.Consumers[runtime.JSONMime] = EasyJSONConsumer()
	return transport
}

func SortedChecks() []Checker {
	pending := map[string]Checker{}
	for _, check := range registeredChecks {
		if _, ok := pending[check.Name]; ok {
			log.Fatal("Found duplicate check:", check.Name)
		}
		pending[check.Name] = check
	}

	result := []Checker{}
	added := map[string]bool{}
	for len(pending) > 0 {
		batch := []Checker{}
		// Found ready tasks
		for _, item := range pending {
			ready := true
			for _, dep := range item.Deps {
				if !added[dep] {
					ready = false
					break
				}
			}
			if ready {
				batch = append(batch, item)
			}
		}
		if len(batch) == 0 {
			log.Fatal("Can't found dependencies for tasks:", pending)
		}
		// Sort batch by name
		sort.Sort(CheckerByName(batch))
		// Add ready tasks to result
		for _, item := range batch {
			added[item.Name] = true
			delete(pending, item.Name)
		}
		result = append(result, batch...)
	}

	return result
}

func templateUid() string {
	return fmt.Sprintf("i%d", atomic.AddInt32(&s_templateUid, 1))
}

func templateAsset(outer, name string) (template.HTML, error) {
	data, err := assets.Asset(name)
	tag := strings.SplitN(outer, " ", 2)[0]
	if err != nil {
		return template.HTML(""), err
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", outer, string(data), tag)), nil
}

func templateDict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func reportTemplate() *template.Template {
	data, err := assets.Asset("template.html")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.
		New("template.html").
		Funcs(template.FuncMap{
			"uid":   templateUid,
			"asset": templateAsset,
			"dict":  templateDict,
		}).
		Parse(string(data))
	if err != nil {
		panic(err)
	}
	return tmpl
}

func Run(url *url.URL, mask *regexp.Regexp, report_file string, keep bool) int {
	total := 0
	failed := 0
	skipped := 0
	broken := map[string]bool{}

	tpl := reportTemplate()
	reports := []*Report{}
	for _, check := range SortedChecks() {
		if (mask != nil) && (mask.FindString(check.Name) == "") {
			continue
		}
		report := Report{
			Checker: check,
		}
		for _, dep := range check.Deps {
			if broken[dep] {
				report.Skip(dep)
			}
		}
		if report.Result != Skipped {
			log.Infof("Run:  %s", check.Name)
			RunCheck(check, &report, url)
		} else {
			log.Noticef("Skip: %s", check.Name)
		}
		total++
		switch report.Result {
		case Skipped:
			broken[check.Name] = true
			skipped++
		case Success:
		default:
			broken[check.Name] = true
			failed++
		}
		reports = append(reports, &report)
		if failed > 0 && !keep {
			break
		}
	}

	if report_file != "" {
		f, err := os.Create(report_file)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()
		err = tpl.Execute(f, struct {
			Total   int
			Failed  int
			Success int
			Skipped int
			Reports []*Report
			Version string
		}{
			Total:   total,
			Failed:  failed,
			Success: total - failed - skipped,
			Skipped: skipped,
			Reports: reports,
			Version: VersionFull(),
		})
		if err != nil {
			panic(err)
		}
	}

	if failed == 0 {
		log.Infof("All tests passed successfully")
	} else {
		skip_info := ""
		if skipped > 0 {
			skip_info = fmt.Sprintf(" (%d skipped)", skipped)
		}
		log.Errorf("Failed %d test of %d%s", failed, total, skip_info)
	}
	return failed
}
