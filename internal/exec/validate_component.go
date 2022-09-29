package exec

import (
	"fmt"
	c "github.com/cloudposse/atmos/pkg/config"
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// ExecuteValidateComponentCmd executes `validate component` command
func ExecuteValidateComponentCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments. The command requires one argument 'componentName'")
	}

	componentName := args[0]

	flags := cmd.Flags()

	stack, err := flags.GetString("stack")
	if err != nil {
		return err
	}

	schemaPath, err := flags.GetString("schema-path")
	if err != nil {
		return err
	}

	schemaType, err := flags.GetString("schema-type")
	if err != nil {
		return err
	}

	_, msg, err := ExecuteValidateComponent(componentName, stack, schemaPath, schemaType)
	if err != nil {
		return err
	}

	u.PrintMessage(msg)
	return nil
}

// ExecuteValidateComponent validates a component in a stack using JsonSchema, OPA or CUE schema documents
func ExecuteValidateComponent(componentName string, stack string, schemaPath string, schemaType string) (bool, string, error) {
	var configAndStacksInfo c.ConfigAndStacksInfo
	configAndStacksInfo.ComponentFromArg = componentName
	configAndStacksInfo.Stack = stack

	configAndStacksInfo.ComponentType = "terraform"
	configAndStacksInfo, err := ProcessStacks(configAndStacksInfo, true)
	if err != nil {
		u.PrintErrorVerbose(err)
		configAndStacksInfo.ComponentType = "helmfile"
		configAndStacksInfo, err = ProcessStacks(configAndStacksInfo, true)
		if err != nil {
			return false, "", err
		}
	}

	componentSection := configAndStacksInfo.ComponentSection

	if schemaPath == "" || schemaType == "" {
		validations, err := FindValidationSection(componentSection)
		if err != nil {
			return false, "", err
		}

		for _, v := range validations {
			schemaPath = v.SchemaPath
			schemaType = v.SchemaType
		}
	}

	return ValidateComponentSection(componentSection, schemaPath, schemaType)
}

// ValidateComponentSection validates the component config using JsonSchema, OPA or CUE schema documents
func ValidateComponentSection(componentSection any, schemaPath string, schemaType string) (bool, string, error) {
	if schemaType != "jsonschema" && schemaType != "opa" && schemaType != "cue" {
		return false, "", fmt.Errorf("invalid schema type '%s'. Supported values: jsonschema, opa, cue", schemaType)
	}

	// Check if the file pointed to by 'schemaPath' exists.
	// If not, join it with the schemas `base_path` from the CLI config
	var filePath string
	if u.FileExists(schemaPath) {
		filePath = schemaPath
	} else {
		switch schemaType {
		case "jsonschema":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.JsonSchema.BasePath, schemaPath)
			}
		case "opa":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.Opa.BasePath, schemaPath)
			}
		case "cue":
			{
				filePath = path.Join(c.Config.BasePath, c.Config.Schemas.Cue.BasePath, schemaPath)
			}
		}

		if !u.FileExists(filePath) {
			return false, "", fmt.Errorf("the file '%s' does not exist for schema type '%s'", schemaPath, schemaType)
		}
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return false, "", err
	}

	schemaText := string(fileContent)

	switch schemaType {
	case "jsonschema":
		{
			return ValidateWithJsonSchema(componentSection, filePath, schemaText)
		}
	case "opa":
		{
			return ValidateWithOpa(componentSection, filePath, schemaText)
		}
	case "cue":
		{
			return ValidateWithCue(componentSection, filePath, schemaText)
		}
	}

	return false, "", fmt.Errorf("invalid 'schema type '%s'. Supported values: jsonschema (default), opa, cue", schemaType)
}

// FindValidationSection finds 'validation' section in the component config
func FindValidationSection(componentSection map[string]any) (c.Validation, error) {
	validationSection := map[any]any{}

	if i, ok := componentSection["settings"].(map[any]any); ok {
		if i2, ok2 := i["validation"].(map[any]any); ok2 {
			validationSection = i2
		}
	}

	var result c.Validation

	err := mapstructure.Decode(validationSection, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
