package cmd

import (
	"fmt"
	"github.com/bmuschko/lets-gopher/templ"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/core"
	"io"
	"os"
	"path"
)

func init() {
	rootCmd.AddCommand(newCreateCmd(rootCmd.OutOrStderr()))
}

type projectCreateCmd struct {
	templateName    string
	templateVersion string
	targetDir       string
	out             io.Writer
	home            templ.Home
}

func newCreateCmd(out io.Writer) *cobra.Command {
	create := &projectCreateCmd{out: out}

	cmd := &cobra.Command{
		Use:   "create [ARGS]",
		Short: "create a new project from a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "the template name", "the template version", "the target project directory"); err != nil {
				return err
			}

			create.templateName = args[0]
			create.templateVersion = args[1]
			create.targetDir = args[2]
			create.home = templ.LetsGopherSettings.Home
			return create.run()
		},
	}

	return cmd
}

func (a *projectCreateCmd) run() error {
	f, err := templ.LoadTemplatesFile(a.home.TemplatesFile())
	if err != nil {
		return err
	}
	if !f.Has(a.templateName) {
		return fmt.Errorf("Template with name %s hasn't been installed", a.templateName)
	}

	templateName := a.templateName + "-" + a.templateVersion
	templateZIP := path.Join(a.home.ArchiveDir(), templateName+".zip")
	archiver := &templ.ZIPArchiver{}

	tb, err := archiver.LoadFile(templateZIP)
	if err != nil {
		return err
	}
	m, err := templ.LoadManifestData(tb)
	if err != nil {
		return err
	}
	r, err := requestParameterValues(m.Parameters)
	if err != nil {
		return err
	}

	err = archiver.Extract(templateZIP, r)
	if err != nil {
		return nil
	}
	err = os.Rename(templateName, a.targetDir)
	return err
}

func requestParameterValues(params []*templ.Parameter) (map[string]string, error) {
	replacements := make(map[string]string)
	if len(params) > 0 {
		core.SetFancyIcons()
	}
	for _, p := range params {
		value := ""
		prompt := &survey.Input{
			Message: "Please enter " + p.Description,
		}
		err := survey.AskOne(prompt, &value, survey.MinLength(1))
		if err != nil {
			return nil, err
		}
		replacements[p.Name] = value
	}

	return replacements, nil
}

//func doGenerateCmd(cmd *cobra.Command, args []string) {
//	config, err := loadConfig()
//	utils.CheckIfError(err)
//
//	availableTemplates := listTemplates()
//	availableTemplateNames := templateNames(availableTemplates)
//
//	if len(availableTemplateNames) > 0 {
//		core.SetFancyIcons()
//		selectedTemplateName := promptTemplate(availableTemplateNames)
//		defaultBasePath := buildDefaultBasePath(config, selectedTemplateName)
//		enteredBasePath := promptBasePath(defaultBasePath)
//
//		templ.GenerateProject(availableTemplates[selectedTemplateName], enteredBasePath)
//	} else {
//		log.Print("No templates found!")
//	}
//}
//
//func listTemplates() map[string]string {
//	var templatesDir = templ.LetsGopherSettings.Home.TemplatesDir()
//	var templates = make(map[string]string)
//
//	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
//		err := utils.CreateDir(templatesDir)
//		utils.CheckIfError(err)
//	} else {
//		files, err := ioutil.ReadDir(templatesDir)
//		utils.CheckIfError(err)
//		for _, f := range files {
//			if f.IsDir() {
//				templates[f.Name()] = filepath.Join(templatesDir, f.Name())
//			}
//		}
//	}
//
//	return templates
//}
//
//
//func promptTemplate(templateNames []string) string {
//	selectedTemplate := ""
//	prompt := &survey.Select{
//		Message: "What template would you like to use to generate a project?",
//		Options: templateNames,
//	}
//	err := survey.AskOne(prompt, &selectedTemplate, nil)
//	utils.CheckIfError(err)
//
//	fmt.Printf("You choose %q\n", selectedTemplate)
//	return selectedTemplate
//}
//
//func promptBasePath(defaultBasePath string) string {
//	enteredBasePath := ""
//	prompt := &survey.Input{
//		Message: "What base path would like to use?",
//		Default: defaultBasePath,
//	}
//	err := survey.AskOne(prompt, &enteredBasePath, survey.MinLength(1))
//	utils.CheckIfError(err)
//
//	fmt.Printf("You choose %q\n", enteredBasePath)
//	return enteredBasePath
//}
//
//func buildDefaultBasePath(genConfig templ.GenConfig, selectedTemplateName string) string {
//	defaultBasePath := ""
//
//	if genConfig.Domain != "" {
//		defaultBasePath += genConfig.Domain
//
//		if !strings.HasSuffix(defaultBasePath, "/") {
//			defaultBasePath += "/"
//		}
//
//		defaultBasePath += selectedTemplateName
//	}
//
//	return defaultBasePath
//}
//
//func templateNames(availableTemplates map[string]string) []string {
//	templateNames := make([]string, 0, len(availableTemplates))
//	for k := range availableTemplates {
//		templateNames = append(templateNames, k)
//	}
//	sort.Strings(templateNames)
//	return templateNames
//}
//
//func loadConfig() (templ.GenConfig, error) {
//	config, err := templ.Load()
//	if err != nil && os.IsNotExist(err) {
//		return templ.GenConfig{}, nil
//	}
//
//	return config, nil
//}
