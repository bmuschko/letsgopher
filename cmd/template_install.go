package cmd

import (
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/bmuschko/letsgopher/template/config"
	"github.com/bmuschko/letsgopher/template/download"
	"github.com/bmuschko/letsgopher/template/environment"
	"github.com/bmuschko/letsgopher/template/storage"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

type templateInstallCmd struct {
	templateURL  string
	templateName string
	out          io.Writer
	home         storage.Home
	downloader   download.Downloader
}

func newTemplateInstallCmd(out io.Writer) *cobra.Command {
	install := &templateInstallCmd{out: out}

	cmd := &cobra.Command{
		Use:   "install [url] [name]",
		Short: "installs a template from a URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the url of the template archive", "name for the template"); err != nil {
				return err
			}

			install.templateURL = args[0]
			install.templateName = args[1]
			install.home = environment.Settings.Home
			install.downloader = &download.TemplateDownloader{Home: environment.Settings.Home, Getter: download.NewHTTPGetter()}
			return install.run()
		},
	}
	return cmd
}

func (c *templateInstallCmd) run() error {
	templateVersion, err := extractTemplateVersion(c.templateURL)
	if err != nil {
		return err
	}
	templateZIP, err := c.downloader.Download(c.templateURL)

	if err != nil {
		return err
	}

	if err := addTemplate(c.templateName, templateVersion, templateZIP, c.home); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%q has been added to your templates\n", c.templateName)
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

func addTemplate(name string, version string, templateZIP string, home storage.Home) error {
	f, err := config.LoadTemplatesFile(home.TemplatesFile())
	if err != nil {
		return err
	}

	if f.Has(name, version) {
		return fmt.Errorf("template with name %q already exists, please specify a different name", name)
	}

	c := config.Template{
		Name:        name,
		Version:     version,
		ArchivePath: templateZIP,
	}
	f.Update(&c)

	return f.WriteFile(home.TemplatesFile(), 0644)
}
