package cmd

import "github.com/spf13/cobra"

func doGenerateCmd(cmd *cobra.Command, args []string) {

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
