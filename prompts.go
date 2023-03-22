package main

import (
	"text/template"
)

var (
	promptDiagnose = template.Must(template.New("diagnose").Parse(`You are a professional kubernetes administrator.
You inspect the object and find out what might cause the error.
Please write down the possible causes in bullet points, using the imperative tense.

The object is
'''
{{ .Data -}}
'''

The cause of the error is
`))
)
