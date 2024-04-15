package global

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

var EmailTemplate *template.Template

func InitTemplate() error {
	log.Println("Reading system email template ...")
	defer log.Println("Reading email template completed.")
	err := readEmailTemplate()
	if err != nil {
		return err
	}
	return nil
}

func readEmailTemplate() error{
	templateBytes, err := os.ReadFile(Config.Server.Mail.Template)
	// log.Println(Config.Server.Mail.Template)
	if err != nil {

		return fmt.Errorf("failed to read template file: %s", err)
	}

	EmailTemplate, err = template.New("email").Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}
	return nil
}
