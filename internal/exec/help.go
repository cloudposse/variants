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
		u.PrintMessage(fmt.Sprintf("'atmos' supports all native '%s' commands.\n", componentType))
		u.PrintMessage("In addition, the 'component' argument and 'stack' flag are required to generate the variables and backend config for the component in the stack.\n")
		u.PrintMessage(fmt.Sprintf("atmos %s <command> <component> -s <stack> [options]", componentType))
		u.PrintMessage(fmt.Sprintf("atmos %s <command> <component> --stack <stack> [options]", componentType))

		if componentType == "terraform" {
			u.PrintMessage("\nAdditions and differences from native terraform:")
			u.PrintMessage(" - before executing other 'terraform' commands, 'atmos' runs 'terraform init'")
			u.PrintMessage(" - you can skip over atmos calling 'terraform init' if you know your project is already in a good working state by using " +
				"the '--skip-init' flag like so 'atmos terraform <command> <component> -s <stack> --skip-init")
			u.PrintMessage(" - 'atmos terraform deploy' command executes 'terraform apply -auto-approve' (sets the '-auto-approve' flag when running 'terraform apply')")
			u.PrintMessage(" - 'atmos terraform deploy' command supports '--deploy-run-init=true/false' flag to enable/disable running 'terraform init' " +
				"before executing the command")
			u.PrintMessage(" - 'atmos terraform apply' and 'atmos terraform deploy' commands support '--from-plan' flag. If the flag is specified, " +
				"the commands will use the planfile previously generated by 'atmos terraform plan' command instead of generating a new planfile")
			u.PrintMessage(" - 'atmos terraform apply' and 'atmos terraform deploy' commands commands support '--planfile' flag to specify the path " +
				"to a planfile. The '--planfile' flag should be used instead of the planfile argument in the native 'terraform apply <planfile>' command")
			u.PrintMessage(" - 'atmos terraform clean' command deletes the '.terraform' folder, '.terraform.lock.hcl' lock file, " +
				"and the previously generated 'planfile', 'varfile' and 'backend.tf.json' file for the specified component and stack. " +
				"Use --skip-lock-file flag to skip deleting the lock file.")
			u.PrintMessage(" - 'atmos terraform workspace' command first runs 'terraform init -reconfigure', then 'terraform workspace select', " +
				"and if the workspace was not created before, it then runs 'terraform workspace new'")
			u.PrintMessage(" - 'atmos terraform import' command searches for 'region' in the variables for the specified component and stack, " +
				"and if it finds it, sets 'AWS_REGION=<region>' ENV var before executing the command")
			u.PrintMessage(" - 'atmos terraform generate backend' command generates a backend config file for an 'atmos' component in a stack")
			u.PrintMessage(" - 'atmos terraform generate backends' command generates backend config files for all 'atmos' components in all stacks")
			u.PrintMessage(" - 'atmos terraform generate varfile' command generates a varfile for an 'atmos' component in a stack")
			u.PrintMessage(" - 'atmos terraform generate varfiles' command generates varfiles for all 'atmos' components in all stacks")
			u.PrintMessage(" - 'atmos terraform shell' command configures an environment for an 'atmos' component in a stack and starts a new shell " +
				"allowing executing all native terraform commands inside the shell without using atmos-specific arguments and flags")
			u.PrintMessage(" - double-dash '--' can be used to signify the end of the options for Atmos and the start of the additional " +
				"native arguments and flags for the 'terraform' commands. " +
				"For example: atmos terraform plan <component> -s <stack> -- -refresh=false -lock=false")

			u.PrintMessage(" - '--append-user-agent' flag sets the TF_APPEND_USER_AGENT environment variable to customize the User-Agent string in Terraform provider requests. " +
				"Example: 'Atmos/<version> (Cloud Posse; +https://atmos.tools)'\n")

		}

		if componentType == "helmfile" {
			u.PrintMessage("\nAdditions and differences from native helmfile:")
			u.PrintMessage(" - 'atmos helmfile generate varfile' command generates a varfile for the component in the stack")
			u.PrintMessage(" - 'atmos helmfile' commands support '[global options]' using the command-line flag '--global-options'. " +
				"Usage: atmos helmfile <command> <component> -s <stack> [command options] [arguments] --global-options=\"--no-color --namespace=test\"")
			u.PrintMessage(" - before executing the 'helmfile' commands, 'atmos' runs 'aws eks update-kubeconfig' to read kubeconfig from " +
				"the EKS cluster and use it to authenticate with the cluster. This can be disabled in 'atmos.yaml' CLI config " +
				"by setting 'components.helmfile.use_eks' to 'false'")
			u.PrintMessage(" - double-dash '--' can be used to signify the end of the options for Atmos and the start of the additional " +
				"native arguments and flags for the 'helmfile' commands")
		}

		fmt.Println()
		err := ExecuteShellCommand(cliConfig, componentType, []string{"--help"}, "", nil, false, "")
		if err != nil {
			return err
		}

	} else if componentType == "terraform" && command == "clean" {
		u.PrintMessage("\n'atmos terraform clean' command deletes the following folders and files from the component's directory:\n\n" +
			" - '.terraform' folder\n" +
			" - folder that the 'TF_DATA_DIR' ENV var points to\n" +
			" - '.terraform.lock.hcl' file\n" +
			" - generated varfile for the component in the stack\n" +
			" - generated planfile for the component in the stack\n" +
			" - generated 'backend.tf.json' file\n\n" +
			"Usage: atmos terraform clean <component> -s <stack> <flags>\n\n" +
			"Use '--skip-lock-file' flag to skip deleting the lock file.\n\n" +
			"For more details refer to https://atmos.tools/cli/commands/terraform/clean\n")
	} else if componentType == "terraform" && command == "deploy" {
		u.PrintMessage("\n'atmos terraform deploy' command executes 'terraform apply -auto-approve' on an Atmos component in an Atmos stack.\n\n" +
			"Usage: atmos terraform deploy <component> -s <stack> <flags>\n\n" +
			"The command automatically sets '-auto-approve' flag when running 'terraform apply'.\n\n" +
			"It supports '--deploy-run-init=true|false' flag to enable/disable running terraform init before executing the command.\n\n" +
			"It supports '--from-plan' flag. If the flag is specified, the command will use the planfile previously generated by 'atmos terraform plan' " +
			"command instead of generating a new planfile.\nNote that in this case, the planfile name is in the format supported by Atmos and is " +
			"saved to the component's folder.\n\n" +
			"It supports '--planfile' flag to specify the path to a planfile.\nThe '--planfile' flag should be used instead of the 'planfile' " +
			"argument in the native 'terraform apply <planfile>' command.\n\n" +
			"For more details refer to https://atmos.tools/cli/commands/terraform/deploy\n")
	} else if componentType == "terraform" && command == "shell" {
		u.PrintMessage("\n'atmos terraform shell' command starts a new SHELL configured with the environment for an Atmos component " +
			"in a Stack to allow executing all native terraform commands\ninside the shell without using the atmos-specific arguments and flags.\n\n" +
			"Usage: atmos terraform shell <component> -s <stack>\n\n" +
			"The command does the following:\n\n" +
			" - Processes the stack config files, generates the required variables for the Atmos component in the stack, and writes them to a file in the component's folder\n" +
			" - Generates a backend config file for the Atmos component in the stack and writes it to a file in the component's folder (or as specified by the Atmos configuration setting)\n" +
			" - Creates a Terraform workspace for the component in the stack\n" +
			" - Drops the user into a separate shell (process) with all the required paths and ENV vars set\n" +
			" - Inside the shell, the user can execute all Terraform commands using the native syntax\n\n" +
			"For more details refer to https://atmos.tools/cli/commands/terraform/shell\n")
	} else if componentType == "terraform" && command == "workspace" {
		u.PrintMessage("\n'atmos terraform workspace' command calculates the Terraform workspace for an Atmos component,\n" +
			"and then executes 'terraform init -reconfigure' and selects the Terraform workspace by executing the 'terraform workspace select' command.\n" +
			"If the workspace does not exist, the command creates it by executing the 'terraform workspace new' command.\n\n" +
			"Usage: atmos terraform workspace <component> -s <stack>\n\n" +
			"For more details refer to https://atmos.tools/cli/commands/terraform/workspace\n")
	} else {
		u.PrintMessage(fmt.Sprintf("'atmos' supports native '%s %s' command with all the options, arguments and flags.\n", componentType, command))
		u.PrintMessage("In addition, 'component' and 'stack' are required in order to generate variables for the component in the stack.\n")
		u.PrintMessage(fmt.Sprintf("atmos %s %s <component> -s <stack> [options]", componentType, command))
		u.PrintMessage(fmt.Sprintf("atmos %s %s <component> --stack <stack> [options]", componentType, command))

		fmt.Println()
		err := ExecuteShellCommand(cliConfig, componentType, []string{command, "--help"}, "", nil, false, "")
		if err != nil {
			return err
		}
	}

	return nil
}
