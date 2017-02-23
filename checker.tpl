<html>
<h1>Report</h1>
{{range .Reports}}
	<h2>[{{ .Result }}] {{ .Checker.Name }}</h2>
	<div>{{ .Checker.Description }}</div>
	{{range .SkippedBy}}
		<div>Skip: {{.}}</div>
	{{end}}
	{{range .Pass}}
		<h3>Pass: {{.Name}}</h3>
		{{range .Messages}}
			<h3>Message: {{ .Delta }}</h3>
		{{end}}
	{{end}}
{{end}}
</html>
