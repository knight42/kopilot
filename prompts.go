package main

import (
	"text/template"
)

var (
	promptDiagnose = template.Must(template.New("diagnose").Parse(`You are a professional kubernetes administrator.
You inspect the object and find out what might cause the error.
If there is no error, please say "Everything is OK".
Response in {{ . Lang -}}
Please write down the possible causes in bullet points, using the imperative tense.

THE OBJECT:
'''
{{ .Data -}}
'''

Remember to write only the most important points and do not write more than a few bullet points.

The cause of the error might be:
`))
)
