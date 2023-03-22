package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func getEnvOrDefault(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

func main() {
	lang := getEnvOrDefault(envKopilotLang, langEN)
	typ := getEnvOrDefault(envKopilotType, typeChatGPT)
	opt := option{
		lang:  lang,
		typ:   typ,
		token: os.Getenv(envKopilotToken),
	}
	cmd := cobra.Command{
		Use: "kopilot",
		Long: fmt.Sprintf(`
You need three ENVs to run Kopilot.
Set %s to specify your token.
Set %s to specify your token type, current type is: %s.
Set %s to specify the language, current language is: %s. Valid options like Chinese, French, Spain, etc.
`, envKopilotToken, envKopilotType, typ, envKopilotLang, lang),
	}
	cmd.AddCommand(newDiagnoseCommand(opt))
	cmd.AddCommand(newAuditCommand(opt))
	_ = cmd.Execute()
}
