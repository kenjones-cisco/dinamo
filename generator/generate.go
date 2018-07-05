package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
)

// Generate file using the specified template file and input data to the specified file location.
func Generate(inputTemplate, outfile string, data map[string]interface{}) error {
	log.WithFields(log.Fields{
		"template": inputTemplate,
		"file":     outfile,
	}).Debugf("generating file with data: %v", data)

	t, err := filepath.Abs(inputTemplate)
	if err != nil {
		return err
	}
	o, err := filepath.Abs(outfile)
	if err != nil {
		return err
	}

	// Load template
	tmpl, err := template.New(filepath.Base(t)).Funcs(sprig.TxtFuncMap()).ParseFiles(t)
	if err != nil {
		return err
	}
	tmpl.Option("missingkey=error")

	f, err := os.Create(o)
	if err != nil {
		return err
	}
	defer f.Close()

	// generate the file
	return tmpl.Execute(f, data)
}

// GenerateUsingEnv generates a file using environment variables as the data for Generate.
func GenerateUsingEnv(inputTemplate, outfile string) error {
	return Generate(inputTemplate, outfile, envMap())
}

func envMap() map[string]interface{} {
	emap := make(map[string]interface{})
	for _, item := range os.Environ() {
		kv := strings.SplitN(item, "=", 2)
		emap[kv[0]] = kv[1]
	}
	return emap
}
