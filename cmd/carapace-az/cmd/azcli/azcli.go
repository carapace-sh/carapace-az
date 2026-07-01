package azcli

import (
	"embed"

	spec "github.com/carapace-sh/carapace-spec"
	"gopkg.in/yaml.v3"
)

//go:embed *.yaml
var f embed.FS

var services map[string]string

func Services() map[string]string {
	return services
}

func Get(name string) (*spec.Command, error) {
	content, err := f.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var command *spec.Command
	if err := yaml.Unmarshal(content, &command); err != nil {
		return nil, err
	}
	return command, nil
}
