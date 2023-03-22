package main

import (
	"text/template"
)

var (
	promptDiagnose = template.Must(template.New("diagnose").Parse(`You are a professional kubernetes administrator.
You inspect the object and find out what might cause the error.
If there is no error, say "Everything is OK".
Write down the possible causes in bullet points, using the imperative tense.
Answer in {{ .Lang }}.

THE OBJECT:
'''
{{ .Data -}}
'''

Remember to write only the most important points and do not write more than a few bullet points.

The cause of the error might be:
`))

	promptAudit = template.Must(template.New("secure").Parse(`You are a professional kubernetes administrator.
You inspect the object and find out what might cause the secure problem and give advise.
If there is no error, say "No security-related issues were found in the provided object.".
Write down the possible causes in bullet points, using the imperative tense.
Answer in {{ .Lang }}.

THE OBJECT:
'''
{{ .Data -}}
'''

Remember to write only the most important points and do not write more than a few bullet points.

The secure problems and related solutions might be:
`))
)
