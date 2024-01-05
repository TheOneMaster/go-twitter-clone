package templates

import (
	"html/template"
	"os"
	"path"
)

const templatePath = "./templates/"

/*
Loads page templates. Already includes components.
*/
func LoadFiles(files ...string) (*template.Template, error) {
	var t *template.Template
	wd, err := os.Getwd()
	if err != nil {
		return t, err
	}
	templatePath := path.Join(wd, "./templates/layouts")

	var filePaths []string
	for _, file := range files {
		filePath := path.Join(templatePath, file)
		filePaths = append(filePaths, filePath)
	}

	t, err = template.ParseFiles(filePaths...)
	if err != nil {
		return t, err
	}

	loadComponents(t)
	return t, err
}

func loadComponents(workingTemplate *template.Template) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	componentsPath := path.Join(wd, templatePath, "./components/*.html")
	workingTemplate.ParseGlob(componentsPath)
}

func LoadFragment(file string) (*template.Template, error) {
	var t *template.Template
	wd, err := os.Getwd()
	if err != nil {
		return t, err
	}
	templatePath := path.Join(wd, "./templates/components", file)

	t, err = template.ParseFiles(templatePath)
	loadComponents(t)
	return t, err
}
