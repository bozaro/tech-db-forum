<html>
<h1>Report</h1>
{{range .Reports}}
	<h2>[{{ .Result }}] {{ .Checker.Name }}</h2>
	<div>{{ .Checker.Description }}</div>
{{end}}
</html>
