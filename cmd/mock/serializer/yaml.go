package serializer

import (
	"io/ioutil"
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
