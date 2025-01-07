package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/pkg/config"
	l "github.com/cloudposse/atmos/pkg/list"
	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

// listComponentsCmd lists atmos components
var listComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Execute 'list components' command",
	Long:  `This command lists all Atmos components or filters components by stacks.`,
	Example: "atmos list components\n" +
		"atmos list components -s <stack>",
	Run: func(cmd *cobra.Command, args []string) {
		// Check Atmos configuration
		checkAtmosConfig()

		stackFlag, _ := cmd.Flags().GetString("stack")
		abstractFlag, _ := cmd.Flags().GetBool("abstract")

		configAndStacksInfo := schema.ConfigAndStacksInfo{}
		atmosConfig, err := config.InitCliConfig(configAndStacksInfo, true)
		if err != nil {
			u.PrintMessageInColor(fmt.Sprintf("Error initializing CLI config: %v", err), color.New(color.FgRed))
			return
		}

		stacksMap, err := e.ExecuteDescribeStacks(atmosConfig, "", nil, nil, nil, false, false, false)
		if err != nil {
			u.PrintMessageInColor(fmt.Sprintf("Error describing stacks: %v", err), color.New(color.FgRed))
			return
		}

		output, err := l.FilterAndListComponents(stackFlag, abstractFlag, stacksMap, atmosConfig.Components.List)
		if err != nil {
			u.PrintMessageInColor(fmt.Sprintf("Error: %v"+"\n", err), color.New(color.FgYellow))
			return
		}

		u.PrintMessageInColor(output, color.New(color.FgGreen))
	},
}

func init() {
	listComponentsCmd.PersistentFlags().StringP("stack", "s", "", "Filter components by stack (e.g., atmos list components -s stack1)")
	listComponentsCmd.PersistentFlags().Bool("abstract", false, "Filter abstract component if true")
	listCmd.AddCommand(listComponentsCmd)
}
