package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	dataFs         = afero.NewOsFs()
	supportedTypes = []string{".json", ".yaml", ".yml"}
)

// DataSources specifies different sources of data
type DataSources struct {
	Data     []string
	DataFile string
	UseEnv   bool
}

// Generate file using the specified template file and input data to the specified file location.
func Generate(inputTemplate, outfile string, sources *DataSources) error {
	data := make(map[string]interface{})

	// handle the sources using the inverse priority (file --> env --> args)
	if sources.DataFile != "" {
		ext, err := fileType(sources.DataFile)
		if err != nil {
			return err
		}
		b, err := readFile(sources.DataFile)
		if err != nil {
			return err
		}
		fileData, err := fileMap(ext, b)
		if err != nil {
			return err
		}
		updateMap(fileData, data)
		log.Debugf("(fileData) data: %v", data)
	}

	if sources.UseEnv {
		envData := listMap(os.Environ())
		updateMap(envData, data)
		log.Debugf("(envData) data: %v", data)
	}

	if len(sources.Data) > 0 {
		argsData := listMap(sources.Data)
		updateMap(argsData, data)
		log.Debugf("(argsData) data: %v", data)
	}

	return generate(inputTemplate, outfile, data)
}

func listMap(list []string) map[string]interface{} {
	log.Debugf("input list: %v", list)
	amap := make(map[string]interface{})
	for _, item := range list {
		kv := strings.SplitN(item, "=", 2)
		amap[kv[0]] = kv[1]
	}
	log.Debugf("map result: %v", amap)
	return amap
}

func fileMap(ext string, data []byte) (map[string]interface{}, error) {
	amap := make(map[string]interface{})

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &amap); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &amap); err != nil {
			return nil, err
		}
	}
	return amap, nil
}

func updateMap(src, dest map[string]interface{}) {
	for k, v := range src {
		if v != nil {
			dest[k] = v
		}
	}
}

func fileType(file string) (string, error) {
	ext := filepath.Ext(file)
	for _, item := range supportedTypes {
		if item == ext {
			return ext, nil
		}
	}
	return "", fmt.Errorf("unsupported file type: %s\nonly %q supported", file, supportedTypes)
}

func readFile(file string) ([]byte, error) {
	f, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	return afero.ReadFile(dataFs, f)
}

func generate(inputTemplate, outfile string, data map[string]interface{}) error {
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

	f, err := dataFs.Create(o)
	if err != nil {
		return err
	}
	defer f.Close()

	// generate the file
	return tmpl.Execute(f, data)
}
