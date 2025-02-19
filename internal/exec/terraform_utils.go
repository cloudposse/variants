package exec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

func checkTerraformConfig(atmosConfig schema.AtmosConfiguration) error {
	if len(atmosConfig.Components.Terraform.BasePath) < 1 {
		return errors.New("Base path to terraform components must be provided in 'components.terraform.base_path' config or " +
			"'ATMOS_COMPONENTS_TERRAFORM_BASE_PATH' ENV variable")
	}

	return nil
}

// cleanTerraformWorkspace deletes the `.terraform/environment` file from the component directory.
// The `.terraform/environment` file contains the name of the currently selected workspace,
// helping Terraform identify the active workspace context for managing your infrastructure.
// We delete the file to prevent the Terraform prompt asking to select the default or the
// previously used workspace. This happens when different backends are used for the same component.
func cleanTerraformWorkspace(atmosConfig schema.AtmosConfiguration, componentPath string) {
	// Get `TF_DATA_DIR` ENV variable, default to `.terraform` if not set
	tfDataDir := os.Getenv("TF_DATA_DIR")
	if tfDataDir == "" {
		tfDataDir = ".terraform"
	}

	// Convert relative path to absolute
	if !filepath.IsAbs(tfDataDir) {
		tfDataDir = filepath.Join(componentPath, tfDataDir)
	}

	// Ensure the path is cleaned properly
	tfDataDir = filepath.Clean(tfDataDir)

	// Construct the full file path
	filePath := filepath.Join(tfDataDir, "environment")

	// Check if the file exists before attempting deletion
	if _, err := os.Stat(filePath); err == nil {
		l.Debug("Terraform environment file found. Proceeding with deletion.", "file", filePath)

		if err := os.Remove(filePath); err != nil {
			l.Debug("Failed to delete Terraform environment file.", "file", filePath, "error", err)
		} else {
			l.Debug("Successfully deleted Terraform environment file.", "file", filePath)
		}
	} else if os.IsNotExist(err) {
		l.Debug("Terraform environment file not found. No action needed.", "file", filePath)
	} else {
		l.Debug("Error checking Terraform environment file.", "file", filePath, "error", err)
	}
}

func shouldProcessStacks(info *schema.ConfigAndStacksInfo) (bool, bool) {
	shouldProcessStacks := true
	shouldCheckStack := true

	if info.SubCommand == "clean" {
		if info.ComponentFromArg == "" {
			shouldProcessStacks = false
		}
		shouldCheckStack = info.Stack != ""

	}

	return shouldProcessStacks, shouldCheckStack
}

func generateBackendConfig(atmosConfig *schema.AtmosConfiguration, info *schema.ConfigAndStacksInfo, workingDir string) error {
	// Auto-generate backend file
	if atmosConfig.Components.Terraform.AutoGenerateBackendFile {
		backendFileName := filepath.Join(workingDir, "backend.tf.json")

		l.Debug("Writing the backend config to file.", "file", backendFileName)

		if !info.DryRun {
			componentBackendConfig, err := generateComponentBackendConfig(info.ComponentBackendType, info.ComponentBackendSection, info.TerraformWorkspace)
			if err != nil {
				return err
			}

			err = u.WriteToFileAsJSON(backendFileName, componentBackendConfig, 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateProviderOverrides(atmosConfig *schema.AtmosConfiguration, info *schema.ConfigAndStacksInfo, workingDir string) error {
	// Generate `providers_override.tf.json` file if the `providers` section is configured
	if len(info.ComponentProvidersSection) > 0 {
		providerOverrideFileName := filepath.Join(workingDir, "providers_override.tf.json")

		l.Debug("Writing the provider overrides to file.", "file", providerOverrideFileName)

		if !info.DryRun {
			providerOverrides := generateComponentProviderOverrides(info.ComponentProvidersSection)
			err := u.WriteToFileAsJSON(providerOverrideFileName, providerOverrides, 0o644)
			return err
		}
	}
	return nil
}

// needProcessTemplatesAndYamlFunctions checks if a Terraform command
// requires the `Go` templates and Atmos YAML functions to be processed
func needProcessTemplatesAndYamlFunctions(command string) bool {
	commandsThatNeedFuncProcessing := []string{
		"init",
		"plan",
		"apply",
		"deploy",
		"destroy",
		"generate",
		"output",
		"clean",
		"shell",
		"write",
		"force-unlock",
		"import",
		"refresh",
		"show",
		"taint",
		"untaint",
		"validate",
		"state list",
		"state mv",
		"state pull",
		"state push",
		"state replace-provider",
		"state rm",
		"state show",
	}
	return u.SliceContainsString(commandsThatNeedFuncProcessing, command)
}

// ExecuteTerraformAffected executes `atmos terraform --affected`
func ExecuteTerraformAffected(cmd *cobra.Command, args []string, info schema.ConfigAndStacksInfo) error {
	// Add these flags here because `atmos describe affected` reads/needs them, but `atmos terraform --affected` does not define them
	cmd.PersistentFlags().String("file", "", "")
	cmd.PersistentFlags().String("format", "yaml", "")
	cmd.PersistentFlags().Bool("verbose", false, "")
	cmd.PersistentFlags().Bool("include-spacelift-admin-stacks", false, "")
	cmd.PersistentFlags().Bool("include-dependents", true, "")
	cmd.PersistentFlags().Bool("include-settings", false, "")
	cmd.PersistentFlags().Bool("upload", false, "")
	cmd.PersistentFlags().Bool("process-templates", true, "")
	cmd.PersistentFlags().Bool("process-functions", true, "")
	cmd.PersistentFlags().StringSlice("skip", nil, "")
	cmd.PersistentFlags().StringP("query", "q", "", "")

	cliArgs, err := parseDescribeAffectedCliArgs(cmd, args)
	if err != nil {
		return err
	}

	cliArgs.IncludeDependents = true
	cliArgs.IncludeSpaceliftAdminStacks = false
	cliArgs.OutputFile = ""
	cliArgs.ProcessTemplates = true
	cliArgs.ProcessYamlFunctions = true
	cliArgs.Skip = nil
	cliArgs.Query = ""

	// https://atmos.tools/cli/commands/describe/affected
	affectedList, _, _, _, err := ExecuteDescribeAffected(cliArgs)
	if err != nil {
		return err
	}

	affectedYaml, err := u.ConvertToYAML(affectedList)
	if err != nil {
		return err
	}
	l.Debug("Affected components:\n" + affectedYaml)

	for _, affected := range affectedList {
		err = executeTerraformAffectedComponent(affected, info)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeTerraformAffectedComponent(affected schema.Affected, info schema.ConfigAndStacksInfo) error {
	// If the affected component is included as dependent in other components, don't process it now,
	// it will be processed in the dependency order
	if !affected.IncludedInDependents {
		info.Component = affected.Component
		info.ComponentFromArg = affected.Component
		info.Stack = affected.Stack

		l.Debug(fmt.Sprintf("Executing: atmos terraform %s %s -s %s", info.SubCommand, affected.Component, affected.Stack))

		err := ExecuteTerraform(info)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExecuteTerraformAll executes `atmos terraform --all`
func ExecuteTerraformAll(cmd *cobra.Command, args []string, info schema.ConfigAndStacksInfo) error {
	return nil
}
