package cmd

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/internal/tui/templates"
	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
	cc "github.com/ivanpirog/coloredcobra"
)

// terraformCmd represents the base command for all terraform sub-commands
var terraformCmd = &cobra.Command{
	Use:                "terraform",
	Aliases:            []string{"tf"},
	Short:              "Execute Terraform commands",
	Long:               `This command executes Terraform commands`,
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		return terraformRun(cmd, cmd, args)
	},
}

// Contains checks if a slice of strings contains an exact match for the target string.
func Contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func terraformRun(cmd *cobra.Command, actualCmd *cobra.Command, args []string) error {
	var argsAfterDoubleDash []string
	var finalArgs = args

	doubleDashIndex := lo.IndexOf(args, "--")
	if doubleDashIndex > 0 {
		finalArgs = lo.Slice(args, 0, doubleDashIndex)
		argsAfterDoubleDash = lo.Slice(args, doubleDashIndex+1, len(args))
	}

	info, _ := e.ProcessCommandLineArgs("terraform", cmd, finalArgs, argsAfterDoubleDash)
	// Exit on help
	if info.NeedHelp {
		if info.SubCommand != "" && info.SubCommand != "--help" && info.SubCommand != "help" {
			suggestions := cmd.SuggestionsFor(args[0])
			if !Contains(suggestions, args[0]) {
				if len(suggestions) > 0 {
					fmt.Printf("Unknown command: '%s'\n\nDid you mean this?\n", args[0])
					for _, suggestion := range suggestions {
						fmt.Printf("  %s\n", suggestion)
					}
				} else {
					fmt.Printf(`Error: Unknkown command %q for %q`+"\n", args[0], cmd.CommandPath())
				}
				fmt.Printf(`Run '%s --help' for usage`+"\n", cmd.CommandPath())
				return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
			}
		}
		// check if this is terraform --help command. TODO: check if this is the best way to do
		if cmd == actualCmd {
			template := templates.GenerateFromBaseTemplate(actualCmd.Use, []templates.HelpTemplateSections{
				templates.LongDescription,
				templates.Usage,
				templates.Aliases,
				templates.Examples,
				templates.AvailableCommands,
				templates.Flags,
				templates.GlobalFlags,
				templates.NativeCommands,
				templates.DoubleDashHelp,
				templates.Footer,
			})
			actualCmd.SetUsageTemplate(template)
			cc.Init(&cc.Config{
				RootCmd:  actualCmd,
				Headings: cc.HiCyan + cc.Bold + cc.Underline,
				Commands: cc.HiGreen + cc.Bold,
				Example:  cc.Italic,
				ExecName: cc.Bold,
				Flags:    cc.Bold,
			})

		}

		actualCmd.Help()
		return nil
	}
	// Check Atmos configuration
	checkAtmosConfig()

	err := e.ExecuteTerraform(info)
	if err != nil {
		u.LogErrorAndExit(schema.AtmosConfiguration{}, err)
	}
	return nil
}

func init() {
	// https://github.com/spf13/cobra/issues/739
	terraformCmd.DisableFlagParsing = true
	terraformCmd.PersistentFlags().StringP("stack", "s", "", "atmos terraform <terraform_command> <component> -s <stack>")
	attachTerraformCommands(terraformCmd)
	RootCmd.AddCommand(terraformCmd)
}
