package audit

import (
	"text/template"
)

type templateData struct {
	Data string
	Lang string
}

var promptAudit = template.Must(template.New("secure").Parse(`You are a professional kubernetes administrator.
You inspect the object and find out the security misconfigurations and give advice.
If there is no problems, say "Everything is OK".
Write down the possible problems in bullet points, using the imperative tense.
Answer in {{ .Lang }}.

THE OBJECT:
'''
{{ .Data -}}
'''

Remember to write only the most important points and do not write more than a few bullet points.

The secure problems and corresponding fix are:
`))
