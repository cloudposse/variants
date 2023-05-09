package exec

import (
	"fmt"
	"reflect"

	"github.com/cloudposse/atmos/pkg/schema"
)

// processConfigSources processes the sources (files) for all sections for a component in a stack
func processConfigSources(
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
) (
	map[string]map[string]any,
	error,
) {
	result := map[string]map[string]any{}

	// `vars` section
	vars := map[string]any{}
	result["vars"] = vars

	for k, v := range configAndStacksInfo.ComponentVarsSection {
		name := k.(string)
		obj := map[string]any{}
		obj["name"] = name
		obj["final_value"] = v
		obj["stack_dependencies"] = processVariableInStacks(configAndStacksInfo, rawStackConfigs, "vars", name)
		vars[name] = obj
	}

	// `env` section
	env := map[string]any{}
	result["env"] = env

	for k, v := range configAndStacksInfo.ComponentEnvSection {
		name := k.(string)
		obj := map[string]any{}
		obj["name"] = name
		obj["final_value"] = v
		obj["stack_dependencies"] = processVariableInStacks(configAndStacksInfo, rawStackConfigs, "env", name)
		env[name] = obj
	}

	// `settings` section
	settings := map[string]any{}
	result["settings"] = settings

	for k, v := range configAndStacksInfo.ComponentSettingsSection {
		name := k.(string)
		obj := map[string]any{}
		obj["name"] = name
		obj["final_value"] = v
		obj["stack_dependencies"] = processVariableInStacks(configAndStacksInfo, rawStackConfigs, "settings", name)
		settings[name] = obj
	}

	return result, nil
}

func processVariableInStacks(
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) []map[string]any {

	result := []map[string]any{}

	// Process the variable for the component in the stack
	// Because we want to show the variable dependencies from higher to lower priority,
	// the order of processing is the reverse order from what Atmos follows when calculating the final variables in the `vars` section
	processComponentVariableInStack(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		configAndStacksInfo,
		rawStackConfigs,
		section,
		variable,
	)

	processComponentVariableInStackImports(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		configAndStacksInfo,
		rawStackConfigs,
		section,
		variable,
	)

	// Process the variable for all the base components in the stack from the inheritance chain
	for _, baseComponent := range configAndStacksInfo.ComponentInheritanceChain {
		processComponentVariableInStack(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			configAndStacksInfo,
			rawStackConfigs,
			section,
			variable,
		)

		processComponentVariableInStackImports(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			configAndStacksInfo,
			rawStackConfigs,
			section,
			variable,
		)
	}

	processComponentTypeVariableInStack(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		configAndStacksInfo,
		rawStackConfigs,
		section,
		variable,
	)

	processGlobalVariableInStack(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		rawStackConfigs,
		section,
		variable,
	)

	for _, baseComponent := range configAndStacksInfo.ComponentInheritanceChain {
		processComponentTypeVariableInStack(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			configAndStacksInfo,
			rawStackConfigs,
			section,
			variable,
		)

		processGlobalVariableInStack(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			rawStackConfigs,
			section,
			variable,
		)
	}

	processComponentTypeVariableInStackImports(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		configAndStacksInfo,
		rawStackConfigs,
		section,
		variable,
	)

	processGlobalVariableInStackImports(
		configAndStacksInfo.ComponentFromArg,
		configAndStacksInfo.StackFile,
		&result,
		rawStackConfigs,
		section,
		variable,
	)

	for _, baseComponent := range configAndStacksInfo.ComponentInheritanceChain {
		processComponentTypeVariableInStackImports(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			configAndStacksInfo,
			rawStackConfigs,
			section,
			variable,
		)

		processGlobalVariableInStackImports(
			baseComponent,
			configAndStacksInfo.StackFile,
			&result,
			rawStackConfigs,
			section,
			variable,
		)
	}

	return result
}

// https://medium.com/swlh/golang-tips-why-pointers-to-slices-are-useful-and-how-ignoring-them-can-lead-to-tricky-bugs-cac90f72e77b
func processComponentVariableInStack(
	component string,
	stackFile string,
	result *[]map[string]any,
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStack, ok := rawStackConfig["stack"]
	if !ok {
		return result
	}

	rawStackMap, ok := rawStack.(map[any]any)
	if !ok {
		return result
	}

	rawStackComponentsSection, ok := rawStackMap["components"]
	if !ok {
		return result
	}

	rawStackComponentsSectionMap, ok := rawStackComponentsSection.(map[any]any)
	if !ok {
		return result
	}

	rawStackComponentTypeSection, ok := rawStackComponentsSectionMap[configAndStacksInfo.ComponentType]
	if !ok {
		return result
	}

	rawStackComponentTypeSectionMap, ok := rawStackComponentTypeSection.(map[any]any)
	if !ok {
		return result
	}

	rawStackComponentSection, ok := rawStackComponentTypeSectionMap[component]
	if !ok {
		return result
	}

	rawStackComponentSectionMap, ok := rawStackComponentSection.(map[any]any)
	if !ok {
		return result
	}

	rawStackVars, ok := rawStackComponentSectionMap[section]
	if !ok {
		return result
	}

	rawStackVarsMap, ok := rawStackVars.(map[any]any)
	if !ok {
		return result
	}

	rawStackVarVal, ok := rawStackVarsMap[variable]
	if !ok {
		return result
	}

	val := map[string]any{
		"stack_file":         stackFile,
		"stack_file_section": fmt.Sprintf("components.%s.%s", configAndStacksInfo.ComponentType, section),
		"variable_value":     rawStackVarVal,
		"dependency_type":    "inline",
	}

	appendSettingDescriptor(result, val)

	return result
}

func processComponentTypeVariableInStack(
	component string,
	stackFile string,
	result *[]map[string]any,
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStack, ok := rawStackConfig["stack"]
	if !ok {
		return result
	}

	rawStackMap, ok := rawStack.(map[any]any)
	if !ok {
		return result
	}

	rawStackComponentTypeSection, ok := rawStackMap[configAndStacksInfo.ComponentType]
	if !ok {
		return result
	}

	rawStackComponentTypeSectionMap, ok := rawStackComponentTypeSection.(map[any]any)
	if !ok {
		return result
	}

	rawStackVars, ok := rawStackComponentTypeSectionMap[section]
	if !ok {
		return result
	}

	rawStackVarsMap, ok := rawStackVars.(map[any]any)
	if !ok {
		return result
	}

	rawStackVarVal, ok := rawStackVarsMap[variable]
	if !ok {
		return result
	}

	val := map[string]any{
		"stack_file":         stackFile,
		"stack_file_section": fmt.Sprintf("%s.%s", configAndStacksInfo.ComponentType, section),
		"variable_value":     rawStackVarVal,
		"dependency_type":    "inline",
	}

	appendSettingDescriptor(result, val)

	return result
}

func processGlobalVariableInStack(
	component string,
	stackFile string,
	result *[]map[string]any,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStack, ok := rawStackConfig["stack"]
	if !ok {
		return result
	}

	rawStackMap, ok := rawStack.(map[any]any)
	if !ok {
		return result
	}

	rawStackVars, ok := rawStackMap[section]
	if !ok {
		return result
	}

	rawStackVarsMap, ok := rawStackVars.(map[any]any)
	if !ok {
		return result
	}

	rawStackVarVal, ok := rawStackVarsMap[variable]
	if !ok {
		return result
	}

	val := map[string]any{
		"stack_file":         stackFile,
		"stack_file_section": section,
		"variable_value":     rawStackVarVal,
		"dependency_type":    "inline",
	}

	appendSettingDescriptor(result, val)

	return result
}

func processComponentVariableInStackImports(
	component string,
	stackFile string,
	result *[]map[string]any,
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStackImports, ok := rawStackConfig["imports"]
	if !ok {
		return result
	}

	rawStackImportsMap, ok := rawStackImports.(map[string]map[any]any)
	if !ok {
		return result
	}

	for impKey, impVal := range rawStackImportsMap {
		rawStackComponentsSection, ok := impVal["components"]
		if !ok {
			continue
		}

		rawStackComponentsSectionMap, ok := rawStackComponentsSection.(map[any]any)
		if !ok {
			continue
		}

		rawStackComponentTypeSection, ok := rawStackComponentsSectionMap[configAndStacksInfo.ComponentType]
		if !ok {
			continue
		}

		rawStackComponentTypeSectionMap, ok := rawStackComponentTypeSection.(map[any]any)
		if !ok {
			continue
		}

		rawStackComponentSection, ok := rawStackComponentTypeSectionMap[component]
		if !ok {
			continue
		}

		rawStackComponentSectionMap, ok := rawStackComponentSection.(map[any]any)
		if !ok {
			continue
		}

		rawStackVars, ok := rawStackComponentSectionMap[section]
		if !ok {
			continue
		}

		rawStackVarsMap, ok := rawStackVars.(map[any]any)
		if !ok {
			continue
		}

		rawStackVarVal, ok := rawStackVarsMap[variable]
		if !ok {
			continue
		}

		val := map[string]any{
			"stack_file":         impKey,
			"stack_file_section": fmt.Sprintf("components.%s.%s", configAndStacksInfo.ComponentType, section),
			"variable_value":     rawStackVarVal,
			"dependency_type":    "import",
		}

		appendSettingDescriptor(result, val)
	}

	return result
}

func processComponentTypeVariableInStackImports(
	component string,
	stackFile string,
	result *[]map[string]any,
	configAndStacksInfo schema.ConfigAndStacksInfo,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStackImports, ok := rawStackConfig["imports"]
	if !ok {
		return result
	}

	rawStackImportsMap, ok := rawStackImports.(map[string]map[any]any)
	if !ok {
		return result
	}

	for impKey, impVal := range rawStackImportsMap {
		rawStackComponentTypeSection, ok := impVal[configAndStacksInfo.ComponentType]
		if !ok {
			continue
		}

		rawStackComponentTypeSectionMap, ok := rawStackComponentTypeSection.(map[any]any)
		if !ok {
			continue
		}

		rawStackVars, ok := rawStackComponentTypeSectionMap[section]
		if !ok {
			continue
		}

		rawStackVarsMap, ok := rawStackVars.(map[any]any)
		if !ok {
			continue
		}

		rawStackVarVal, ok := rawStackVarsMap[variable]
		if !ok {
			continue
		}

		val := map[string]any{
			"stack_file":         impKey,
			"stack_file_section": fmt.Sprintf("%s.%s", configAndStacksInfo.ComponentType, section),
			"variable_value":     rawStackVarVal,
			"dependency_type":    "import",
		}

		appendSettingDescriptor(result, val)
	}

	return result
}

func processGlobalVariableInStackImports(
	component string,
	stackFile string,
	result *[]map[string]any,
	rawStackConfigs map[string]map[string]any,
	section string,
	variable string,
) *[]map[string]any {

	rawStackConfig, ok := rawStackConfigs[stackFile]
	if !ok {
		return result
	}

	rawStackImports, ok := rawStackConfig["imports"]
	if !ok {
		return result
	}

	rawStackImportsMap, ok := rawStackImports.(map[string]map[any]any)
	if !ok {
		return result
	}

	for impKey, impVal := range rawStackImportsMap {
		rawStackVars, ok := impVal[section]
		if !ok {
			continue
		}

		rawStackVarsMap, ok := rawStackVars.(map[any]any)
		if !ok {
			continue
		}

		rawStackVarVal, ok := rawStackVarsMap[variable]
		if !ok {
			continue
		}

		val := map[string]any{
			"stack_file":         impKey,
			"stack_file_section": section,
			"variable_value":     rawStackVarVal,
			"dependency_type":    "import",
		}

		appendSettingDescriptor(result, val)
	}

	return result
}

func appendSettingDescriptor(result *[]map[string]any, descriptor map[string]any) {
	for _, item := range *result {
		if reflect.DeepEqual(item, descriptor) {
			return
		}
	}
	*result = append(*result, descriptor)
}
