package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/go-openapi/runtime"
	http_transport "github.com/go-openapi/runtime/client"
	"log"
	"net/http"
	"net/url"
)

type Checker struct {
	// Имя текущей проверки.
	Name string
	// Описание текущей проверки.
	Description string
	// Функция для текущей проверки.
	FnCheck func(c *client.Forum)
	// Тесты, без которых проверка не имеет смысл.
	Deps []string
}

type CheckerTransport struct {
	t      runtime.ClientTransport
	report *Report
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := NewValidator(operation.Context, self.report)
	operation.Client = &http.Client{Transport: tracker}
	return self.t.Submit(operation)
}

func Checkpoint(c *client.Forum, message string) bool {
	return c.Transport.(*CheckerTransport).report.Checkpoint(message)
}

var checks []Checker

func Register(checker Checker) {
	checks = append(checks, checker)
}

func RunCheck(check Checker, report *Report, cfg *client.TransportConfig) {
	report.result = REPORT_SUCCESS
	transport := http_transport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	defer func() {
		if r := recover(); r != nil {
			report.AddError(r)
		}
	}()
	check.FnCheck(client.New(&CheckerTransport{transport, report}, nil))
}

func Run(url *url.URL) int {
	total := 0
	failed := 0
	skipped := 0

	cfg := client.DefaultTransportConfig().WithHost(url.Host).WithSchemes([]string{url.Scheme}).WithBasePath(url.Path)
	for _, check := range checks {
		log.Printf("=== RUN:  %s", check.Name)
		report := Report{}
		RunCheck(check, &report, cfg)
		if report.result != REPORT_SUCCESS {
			report.Show()
		}
		var result string
		total++
		switch report.result {
		case REPORT_SKIPPED:
			skipped++
			result = "SKIPPED"
		case REPORT_SUCCESS:
			result = "OK"
		default:
			failed++
			result = "FAILED"
		}
		log.Printf("--- DONE: %s (%s)", check.Name, result)
	}
	log.Printf("RESULT: %d total, %d skipped, %d failed)", total, skipped, failed)
	return failed
}
