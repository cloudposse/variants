package exec

import (
	"fmt"
	c "github.com/cloudposse/atmos/pkg/config"
	s "github.com/cloudposse/atmos/pkg/stack"
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

// ExecuteValidateStacks executes `validate stacks` command
func ExecuteValidateStacks(cmd *cobra.Command, args []string) error {
	cliConfig, err := c.InitCliConfig(c.ConfigAndStacksInfo{}, true)
	if err != nil {
		return err
	}

	// Include (process and validate) all YAML files in the `stacks` folder in all subfolders
	includedPaths := []string{"**/*"}
	// Don't exclude any YAML files for validation
	excludedPaths := []string{}
	includeStackAbsPaths, err := u.JoinAbsolutePathWithPaths(cliConfig.StacksBaseAbsolutePath, includedPaths)
	if err != nil {
		return err
	}

	stackConfigFilesAbsolutePaths, _, err := c.FindAllStackConfigsInPaths(cliConfig, includeStackAbsPaths, excludedPaths)
	if err != nil {
		return err
	}

	u.PrintInfo(fmt.Sprintf("Validating all YAML files in the '%s' folder and all subfolders\n",
		path.Join(cliConfig.BasePath, cliConfig.Stacks.BasePath)))

	var errorMessages []string
	for _, filePath := range stackConfigFilesAbsolutePaths {
		stackConfig, importsConfig, err := s.ProcessYAMLConfigFile(cliConfig.StacksBaseAbsolutePath, filePath, map[string]map[any]any{})
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}

		componentStackMap := map[string]map[string][]string{}
		_, err = s.ProcessStackConfig(
			cliConfig.StacksBaseAbsolutePath,
			cliConfig.TerraformDirAbsolutePath,
			cliConfig.HelmfileDirAbsolutePath,
			filePath,
			stackConfig,
			false,
			true,
			"",
			componentStackMap,
			importsConfig,
			false)
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, "\n\n"))
	}

	return nil
}
