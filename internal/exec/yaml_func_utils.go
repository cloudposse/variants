package exec

import (
	"fmt"
	"strings"

	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

func ProcessCustomYamlTags(
	cliConfig schema.CliConfiguration,
	input schema.AtmosSectionMapType,
	currentStack string,
) (schema.AtmosSectionMapType, error) {
	return processNodes(cliConfig, input, currentStack), nil
}

func processNodes(
	cliConfig schema.CliConfiguration,
	data map[string]any,
	currentStack string,
) map[string]any {
	newMap := make(map[string]any)
	var recurse func(any) any

	recurse = func(node any) any {
		switch v := node.(type) {
		case string:
			return processCustomTags(cliConfig, v, currentStack)

		case map[string]any:
			newNestedMap := make(map[string]any)
			for k, val := range v {
				newNestedMap[k] = recurse(val)
			}
			return newNestedMap

		case []any:
			newSlice := make([]any, len(v))
			for i, val := range v {
				newSlice[i] = recurse(val)
			}
			return newSlice

		default:
			return v
		}
	}

	for k, v := range data {
		newMap[k] = recurse(v)
	}

	return newMap
}

func processCustomTags(
	cliConfig schema.CliConfiguration,
	input string,
	currentStack string,
) any {

	switch {
	case strings.HasPrefix(input, u.AtmosYamlFuncTemplate):
		return processTagTemplate(cliConfig, input, currentStack)
	case strings.HasPrefix(input, u.AtmosYamlFuncExec):
		return processTagExec(cliConfig, input, currentStack)
	case strings.HasPrefix(input, u.AtmosYamlFuncStore):
		return processTagStore(cliConfig, input, currentStack)
	case strings.HasPrefix(input, u.AtmosYamlFuncTerraformOutput):
		return processTagTerraformOutput(cliConfig, input, currentStack)
	default:
		// If any other YAML explicit type (not currently supported by Atmos) is used, return it w/o processing
		return input
	}
}

func getStringAfterTag(input string, tag string) (string, error) {
	str := strings.TrimPrefix(input, tag)
	str = strings.TrimSpace(str)

	if str == "" {
		err := fmt.Errorf("invalid Atmos YAML function: %s", input)
		return "", err
	}

	return str, nil
}
