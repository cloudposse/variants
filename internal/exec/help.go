package exec

import (
	"fmt"

	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

// processHelp processes help commands
func processHelp(componentType string, command string) error {
	cliConfig := schema.CliConfiguration{}
	if len(command) == 0 {
		u.LogMessage(cliConfig, fmt.Sprintf("'atmos' supports all native '%s' commands.\n", componentType))
		u.LogMessage(cliConfig, "In addition, the 'component' argument and 'stack' flag are required to generate the variables and backend config for the component in the stack.\n")
		u.LogInfo(cliConfig, fmt.Sprintf("atmos %s <command> <component> -s <stack> [options]", componentType))
		u.LogInfo(cliConfig, fmt.Sprintf("atmos %s <command> <component> --stack <stack> [options]", componentType))

		if componentType == "terraform" {
			u.LogInfo(cliConfig, "\nAdditions and differences from native terraform:")
			u.LogMessage(cliConfig, " - before executing other 'terraform' commands, 'atmos' runs 'terraform init'")
			u.LogMessage(cliConfig, " - you can skip over atmos calling 'terraform init' if you know your project is already in a good working state by using "+
				"the '--skip-init' flag like so 'atmos terraform <command> <component> -s <stack> --skip-init")
			u.LogMessage(cliConfig, " - 'atmos terraform deploy' command executes 'terraform apply -auto-approve' (sets the '-auto-approve' flag when running 'terraform apply')")
			u.LogMessage(cliConfig, " - 'atmos terraform deploy' command supports '--deploy-run-init=true/false' flag to enable/disable running 'terraform init' "+
				"before executing the command")
			u.LogMessage(cliConfig, " - 'atmos terraform apply' and 'atmos terraform deploy' commands support '--from-plan' flag. If the flag is specified, "+
				"the commands will use the planfile previously generated by 'atmos terraform plan' command instead of generating a new planfile")
			u.LogMessage(cliConfig, " - 'atmos terraform apply' and 'atmos terraform deploy' commands commands support '--planfile' flag to specify the path "+
				"to a planfile. The '--planfile' flag should be used instead of the planfile argument in the native 'terraform apply <planfile>' command")
			u.LogMessage(cliConfig, " - 'atmos terraform clean' command deletes the '.terraform' folder, '.terraform.lock.hcl' lock file, "+
				"and the previously generated 'planfile' and 'varfile' for the specified component and stack")
			u.LogMessage(cliConfig, " - 'atmos terraform workspace' command first runs 'terraform init -reconfigure', then 'terraform workspace select', "+
				"and if the workspace was not created before, it then runs 'terraform workspace new'")
			u.LogMessage(cliConfig, " - 'atmos terraform import' command searches for 'region' in the variables for the specified component and stack, "+
				"and if it finds it, sets 'AWS_REGION=<region>' ENV var before executing the command")
			u.LogMessage(cliConfig, " - 'atmos terraform generate backend' command generates a backend config file for an 'atmos' component in a stack")
			u.LogMessage(cliConfig, " - 'atmos terraform generate backends' command generates backend config files for all 'atmos' components in all stacks")
			u.LogMessage(cliConfig, " - 'atmos terraform generate varfile' command generates a varfile for an 'atmos' component in a stack")
			u.LogMessage(cliConfig, " - 'atmos terraform generate varfiles' command generates varfiles for all 'atmos' components in all stacks")
			u.LogMessage(cliConfig, " - 'atmos terraform shell' command configures an environment for an 'atmos' component in a stack and starts a new shell "+
				"allowing executing all native terraform commands inside the shell without using atmos-specific arguments and flags")
		}

		if componentType == "helmfile" {
			u.LogInfo(cliConfig, "\nAdditions and differences from native helmfile:")
			u.LogMessage(cliConfig, " - 'atmos helmfile generate varfile' command generates a varfile for the component in the stack")
			u.LogMessage(cliConfig, " - 'atmos helmfile' commands support '[global options]' using the command-line flag '--global-options'. "+
				"Usage: atmos helmfile <command> <component> -s <stack> [command options] [arguments] --global-options=\"--no-color --namespace=test\"")
			u.LogMessage(cliConfig, " - before executing the 'helmfile' commands, 'atmos' runs 'aws eks update-kubeconfig' to read kubeconfig from "+
				"the EKS cluster and use it to authenticate with the cluster. This can be disabled in 'atmos.yaml' CLI config "+
				"by setting 'components.helmfile.use_eks' to 'false'")
		}

		err := ExecuteShellCommand(cliConfig, componentType, []string{"--help"}, "", nil, false, true, "")
		if err != nil {
			return err
		}
	} else {
		u.LogMessage(cliConfig, fmt.Sprintf("'atmos' supports native '%s %s' command with all the options, arguments and flags.\n", componentType, command))
		u.LogMessage(cliConfig, "In addition, 'component' and 'stack' are required in order to generate variables for the component in the stack.\n")
		u.LogInfo(cliConfig, fmt.Sprintf("atmos %s %s <component> -s <stack> [options]", componentType, command))
		u.LogInfo(cliConfig, fmt.Sprintf("atmos %s %s <component> --stack <stack> [options]", componentType, command))

		err := ExecuteShellCommand(cliConfig, componentType, []string{command, "--help"}, "", nil, false, true, "")
		if err != nil {
			return err
		}
	}

	return nil
}
