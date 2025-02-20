package stack

import (
	"github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/pkg/schema"
)

// and returns a list of stack configs.
func ProcessYAMLConfigFiles(
	atmosConfig schema.AtmosConfiguration,
	stacksBasePath string,
	terraformComponentsBasePath string,
	helmfileComponentsBasePath string,
	filePaths []string,
	processStackDeps bool,
	processComponentDeps bool,
	ignoreMissingFiles bool,
) (
	[]string,
	map[string]any,
	map[string]map[string]any,
	error,
) {
	return exec.ProcessYAMLConfigFiles(
		atmosConfig,
		stacksBasePath,
		terraformComponentsBasePath,
		helmfileComponentsBasePath,
		filePaths,
		processStackDeps,
		processComponentDeps,
		ignoreMissingFiles,
	)
}

func ProcessYAMLConfigFile(
	atmosConfig schema.AtmosConfiguration,
	basePath string,
	filePath string,
	importsConfig map[string]map[string]any,
	context map[string]any,
	ignoreMissingFiles bool,
	skipTemplatesProcessingInImports bool,
	ignoreMissingTemplateValues bool,
	skipIfMissing bool,
	parentTerraformOverrides map[string]any,
	parentHelmfileOverrides map[string]any,
	atmosManifestJsonSchemaFilePath string,
) (
	map[string]any,
	map[string]map[string]any,
	map[string]any,
	map[string]any,
	map[string]any,
	error,
) {
	return exec.ProcessYAMLConfigFile(
		atmosConfig,
		basePath,
		filePath,
		importsConfig,
		context,
		ignoreMissingFiles,
		skipTemplatesProcessingInImports,
		ignoreMissingTemplateValues,
		skipIfMissing,
		parentTerraformOverrides,
		parentHelmfileOverrides,
		atmosManifestJsonSchemaFilePath,
	)
}

// and returns the final stack configuration for all Terraform and helmfile components.
func ProcessStackConfig(
	atmosConfig schema.AtmosConfiguration,
	stacksBasePath string,
	terraformComponentsBasePath string,
	helmfileComponentsBasePath string,
	stack string,
	config map[string]any,
	processStackDeps bool,
	processComponentDeps bool,
	componentTypeFilter string,
	componentStackMap map[string]map[string][]string,
	importsConfig map[string]map[string]any,
	checkBaseComponentExists bool,
) (map[string]any, error) {
	return exec.ProcessStackConfig(
		atmosConfig,
		stacksBasePath,
		terraformComponentsBasePath,
		helmfileComponentsBasePath,
		stack,
		config,
		processStackDeps,
		processComponentDeps,
		componentTypeFilter,
		componentStackMap,
		importsConfig,
		checkBaseComponentExists,
	)
}
