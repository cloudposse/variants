package cmd

import _ "embed"

//go:embed markdown/atmos_terraform_usage.md
var terraform string

//go:embed markdown/atmos_terraform_plan_usage.md
var terraformPlan string

//go:embed markdown/atmos_terraform_apply_usage.md
var terraformApply string

//go:embed markdown/atmos_workflow_usage.md
var workflow string

//go:embed markdown/atmos_about_usage.md
var atmosAbout string

type ExampleContent struct {
	Content    string
	Suggestion string
}

var doubleDashHint string = "Use double dashes to separate Atmos-specific options from native arguments and flags for the command."

var examples map[string]ExampleContent = map[string]ExampleContent{
	"atmos_terraform": {
		Content:    terraform,
		Suggestion: "https://atmos.tools/cli/commands/terraform/usage",
	},
	"atmos_terraform_plan": {
		Content: terraformPlan,
		// TODO: We should update this once we have a page for terraform plan
		Suggestion: "https://atmos.tools/cli/commands/terraform/usage",
	},
	"atmos_terraform_apply": {
		Content: terraformApply,
		// TODO: We should update this once we have a page for terraform plan
		Suggestion: "https://atmos.tools/cli/commands/terraform/usage",
	},
	"atmos_workflow": {
		Content:    workflow,
		Suggestion: "https://atmos.tools/cli/commands/workflow/",
	},
	"atmos_about": {
		Content: atmosAbout,
	},
}
