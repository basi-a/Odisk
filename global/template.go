package global

import (
	"html/template"
	"log"
	"os"
)

var EmailTemplate *template.Template

func InitTemplate() {
	readEmailTemplate()
}

func readEmailTemplate() {
	templateBytes, err := os.ReadFile(Config.Server.Mail.Template)
	// log.Println(Config.Server.Mail.Template)
	if err != nil {
		log.Fatalf("Failed to read template file: %v", err)
	}

	EmailTemplate, err = template.New("email").Parse(string(templateBytes))
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
}
