package diagnose

import (
	"text/template"
)

type templateData struct {
	Data string
	Lang string
}

var promptDiagnose = template.Must(template.New("diagnose").Parse(`You are a professional kubernetes administrator.
Carefully read the provided information, being certain to spell out the diagnosis & reasoning, and don't skip any steps.
Answer in {{ .Lang }}.

---

{{ .Data }}

---

What is wrong with this object and how to fix it?
`))
