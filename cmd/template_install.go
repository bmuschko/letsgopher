package cmd

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/bmuschko/lets-gopher/templ/config"
	"github.com/bmuschko/lets-gopher/templ/download"
	"github.com/bmuschko/lets-gopher/templ/environment"
	"github.com/bmuschko/lets-gopher/templ/path"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

type templateInstallCmd struct {
	templateURL  string
	templateName string
	out          io.Writer
	home         path.Home
}

func newTemplateInstallCmd(out io.Writer) *cobra.Command {
	add := &templateInstallCmd{out: out}

	cmd := &cobra.Command{
		Use:   "install [url] [name]",
		Short: "installs a template from a URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the url of the template archive", "name for the template"); err != nil {
				return err
			}

			add.templateURL = args[0]
			add.templateName = args[1]
			add.home = environment.Settings.Home
			return add.run()
		},
	}
	return cmd
}

func (a *templateInstallCmd) run() error {
	templateVersion, err := extractTemplateVersion(a.templateURL)
	if err != nil {
		return err
	}
	downloader := &download.TemplateDownloader{Home: environment.Settings.Home, Getter: download.NewHTTPGetter()}
	templateZIP, err := downloader.Download(a.templateURL)

	if err != nil {
		return nil
	}

	if err := addTemplate(a.templateName, templateVersion, templateZIP, a.home); err != nil {
		return err
	}
	fmt.Fprintf(a.out, "%q has been added to your templates\n", a.templateName)
	return nil
}

func extractTemplateVersion(url string) (string, error) {
	lastSlash := strings.LastIndex(url, "/")
	lastDot := strings.LastIndex(url, ".")
	r := []rune(url)
	templateName := string(r[lastSlash+1 : lastDot])
	versionSeparatorIndex := strings.LastIndex(templateName, "-")

	if versionSeparatorIndex == -1 {
		return "", errors.New("template archive file name needs to contain a version separated by a dash character")
	}

	t := []rune(templateName)
	templateVersion := string(t[versionSeparatorIndex+1 : len(templateName)])
	parsedVersion, err := semver.Make(templateVersion)
	if err != nil {
		return "", err
	}
	return parsedVersion.String(), nil
}

func addTemplate(name string, version string, templateZIP string, home path.Home) error {
	f, err := config.LoadTemplatesFile(home.TemplatesFile())
	if err != nil {
		return err
	}

	if f.Has(name, version) {
		return fmt.Errorf("template with name (%s) already exists, please specify a different name", name)
	}

	c := config.Template{
		Name:        name,
		Version:     version,
		ArchivePath: templateZIP,
	}
	f.Update(&c)

	return f.WriteFile(home.TemplatesFile(), 0644)
}
