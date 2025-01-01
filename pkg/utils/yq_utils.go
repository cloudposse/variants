// https://github.com/mikefarah/yq
// https://mikefarah.gitbook.io/yq
// https://mikefarah.gitbook.io/yq/recipes
// https://mikefarah.gitbook.io/yq/operators/pipe

package utils

import (
	"fmt"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func EvaluateYqExpression(data map[string]any, expression string) (any, error) {
	evaluator := yqlib.NewStringEvaluator()

	yaml, err := ConvertToYAML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert data to YAML: %w", err)
	}

	pref := yqlib.YamlPreferences{
		Indent:                      2,
		ColorsEnabled:               false,
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          true,
		UnwrapScalar:                true,
		EvaluateTogether:            false,
	}

	encoder := yqlib.NewYamlEncoder(pref)
	decoder := yqlib.NewYamlDecoder(pref)

	result, err := evaluator.Evaluate(expression, yaml, encoder, decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate YQ expression '%s': %w", expression, err)
	}

	res, err := UnmarshalYAML[any](result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert data to YAML: %w", err)
	}

	return res, nil
}
