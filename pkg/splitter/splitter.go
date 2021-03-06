package splitter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	yaml "github.com/replicatedhq/yaml/v3"
)

func SplitYAML(input []byte) (map[string][]byte, error) {
	outputFiles := map[string][]byte{}
	docs := bytes.Split(input, []byte("\n---\n"))

	for _, doc := range docs {
		if bytes.HasPrefix(doc, []byte("---\n")) {
			doc = doc[4:]
		}

		if len(doc) == 0 {
			continue
		}

		filename, err := generateName(doc)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate name")
		}

		outputFiles[filename] = doc
	}

	return outputFiles, nil
}

func generateName(content []byte) (string, error) {
	o := OverlySimpleGVK{}

	if err := yaml.Unmarshal(content, &o); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal yaml")
	}

	return fmt.Sprintf("%s-%s.yaml", o.Metadata.Name, strings.ToLower(o.Kind)), nil
}
