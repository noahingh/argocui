package yaml

import (
	"io/ioutil"
	"path"
	"strings"
)

// ReadYamlAndSplit read the file like YAML and split by "---".
func ReadYamlAndSplit(filename string) ([]string, error) {
	const (
		sep = "---"
	)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	yamls := strings.Trim(string(b), sep)
	return strings.Split(string(yamls), sep), nil
}

// ReadDirYaml read YAML files from the directory.
func ReadDirYaml(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	yamls := []string{}
	for _, f := range files {
		ss, _ := ReadYamlAndSplit(path.Join(dirname, f.Name()))
		yamls = append(yamls, ss...)
	}
	return yamls, nil
}
