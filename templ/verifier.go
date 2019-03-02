package templ

import (
	"fmt"
	"path"
	"strings"
)

type TemplateURLVerifier struct {
	URL string
}

func (t TemplateURLVerifier) Verify() error {
	err := validateExt(t.URL)

	if err != nil {
		return err
	}

	err = validateVersion(t.URL)

	if err != nil {
		return err
	}

	return nil
}

func validateExt(url string) error {
	ext := path.Ext(url)

	if ext != "" || ext != ".zip" {
		return fmt.Errorf("URL  %s needs to point to ZIP file: %s", url)
	}

	return nil
}

func validateVersion(url string) error {
	lastDotSlash := strings.LastIndex(url, "/")
	lastDotIndex := strings.LastIndex(url, ".")
	r := []rune(url)
	fullTemplateName := string(r[lastDotSlash:lastDotIndex])

	if !strings.Contains(fullTemplateName, "-") {
		return fmt.Errorf("Template %s does not contain hypen character to separate name from version", fullTemplateName)
	}

	return nil
}