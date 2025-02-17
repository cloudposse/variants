package cmd

import (
	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	u "github.com/cloudposse/atmos/pkg/utils"
)

// vendorPullCmd executes 'vendor pull' CLI commands
var vendorPullCmd = &cobra.Command{
	Use:                "pull",
	Short:              "Pull the latest vendor configurations or dependencies",
	Long:               "Pull and update vendor-specific configurations or dependencies to ensure the project has the latest required resources.",
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: false},
	Args:               cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// WithStackValidation is a functional option that enables/disables stack configuration validation
		// based on whether the --stack flag is provided
		checkAtmosConfig(WithStackValidation(cmd.Flag("stack").Changed))

		err := e.ExecuteVendorPullCmd(cmd, args)
		if err != nil {
			u.PrintErrorMarkdownAndExit("", err, "")
		}
	},
}

func init() {
	vendorPullCmd.PersistentFlags().StringP("component", "c", "", "Only vendor the specified component: atmos vendor pull --component &ltcomponent&gt")
	vendorPullCmd.RegisterFlagCompletionFunc("component", ComponentsArgCompletion)
	vendorPullCmd.PersistentFlags().StringP("stack", "s", "", "Only vendor the specified stack: atmos vendor pull --stack &ltstack&gt")
	AddStackCompletion(vendorPullCmd)
	vendorPullCmd.PersistentFlags().StringP("type", "t", "terraform", "atmos vendor pull --component &ltcomponent&gt --type=terraform|helmfile")
	vendorPullCmd.PersistentFlags().Bool("dry-run", false, "atmos vendor pull --component &ltcomponent&gt --dry-run")
	vendorPullCmd.PersistentFlags().String("tags", "", "Only vendor the components that have the specified tags: atmos vendor pull --tags=dev,test")
	vendorPullCmd.PersistentFlags().Bool("everything", false, "Vendor all components: atmos vendor pull --everything")
	vendorCmd.AddCommand(vendorPullCmd)
}
