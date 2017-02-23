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
			<h3>URL: {{ .Url }}</h3>			
			<h4>Request</h4>
			{{with .Request}}
			{{.Title}}
			<ul>
			{{range $key, $value := .Header}}
				{{range $value}}
				<li>
				{{ $key }}: {{ . }}
				</li>
				{{end}}
			{{end}}
			</ul>
			<pre>{{.Body}}</pre>
			{{end}}

			{{if ne .Delta "" }}
			<h4>Delta</h4>
			<table>
			{{ .Delta }}
			</table>
			{{end}}

			<h4>Response</h4>
			{{with .Response}}
			{{.Title}}
			<ul>
			{{range $key, $value := .Header}}
				{{range $value}}
				<li>
				{{ $key }}: {{ . }}
				</li>
				{{end}}
			{{end}}
			</ul>
			<pre>{{.Body}}</pre>
			{{end}}

			<h4>Example</h4>
			{{with .Example}}
			{{.Title}}
			<ul>
			{{range $key, $value := .Header}}
				{{range $value}}
				<li>
				{{ $key }}: {{ . }}
				</li>
				{{end}}
			{{end}}
			</ul>
			<pre>{{.Body}}</pre>
			{{end}}

		{{end}}
	{{end}}
{{end}}
</html>
