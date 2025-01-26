package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/pkg/schema"
	"github.com/cloudposse/atmos/pkg/ui/markdown"
	u "github.com/cloudposse/atmos/pkg/utils"
)

//go:embed markdown/workflow.md
var workflowMarkdown string

// ErrorMessage represents a structured error message
type ErrorMessage struct {
	Title      string
	Details    string
	Suggestion string
}

// renderError renders an error message using the markdown renderer
func renderError(msg ErrorMessage) error {
	renderer, err := markdown.NewTerminalMarkdownRenderer(atmosConfig)
	if err != nil {
		return err
	}
	rendered, err := renderer.RenderError(msg.Title, msg.Details, msg.Suggestion)
	if err != nil {
		return fmt.Errorf("failed to render error message: %w", err)
	}

	fmt.Print(rendered + "\n")
	return nil
}

// workflowCmd executes a workflow
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Run predefined tasks using workflows",
	Long:  `Run predefined workflows as an alternative to traditional task runners. Workflows enable you to automate and manage infrastructure and operational tasks specified in configuration files.`,
	Example: "atmos workflow\n" +
		"atmos workflow <name> --file <file>\n" +
		"atmos workflow <name> --file <file> --stack <stack>\n" +
		"atmos workflow <name> --file <file> --from-step <step-name>\n\n" +
		"To resume the workflow from this step, run:\n" +
		"atmos workflow deploy-infra --file workflow1 --from-step deploy-vpc\n\n" +
		"For more details refer to https://atmos.tools/cli/commands/workflow/",
	FParseErrWhitelist: struct{ UnknownFlags bool }{UnknownFlags: false},
	Run: func(cmd *cobra.Command, args []string) {
		handleHelpRequest(cmd, args)
		// If no arguments are provided, start the workflow UI
		if len(args) == 0 {
			err := e.ExecuteWorkflowCmd(cmd, args)
			if err != nil {
				u.LogErrorAndExit(schema.AtmosConfiguration{}, err)
			}
			return
		}

		// Get the --file flag value
		workflowFile, _ := cmd.Flags().GetString("file")

		// If no file is provided, show invalid command error with usage information
		if workflowFile == "" {
			cmd.Usage()
		}

		// Execute the workflow command
		err := e.ExecuteWorkflowCmd(cmd, args)
		if err != nil {
			// Format common error messages
			if strings.Contains(err.Error(), "does not exist") {
				err := renderError(ErrorMessage{
					Title:      "Workflow File Not Found",
					Details:    fmt.Sprintf("`%v` was not found", workflowFile),
					Suggestion: "",
				})
				if err != nil {
					u.LogErrorAndExit(schema.AtmosConfiguration{}, err)
				}
			} else if strings.Contains(err.Error(), "does not have the") {
				err := renderError(ErrorMessage{
					Title:      "Invalid Workflow",
					Details:    err.Error(),
					Suggestion: "",
				})
				if err != nil {
					u.LogErrorAndExit(schema.AtmosConfiguration{}, err)
				}
			} else {
				// For other errors, use the standard error handler
				u.LogErrorAndExit(schema.AtmosConfiguration{}, err)
			}
			os.Exit(1)
		}
	},
}

func init() {
	workflowCmd.DisableFlagParsing = false
	workflowCmd.PersistentFlags().StringP("file", "f", "", "atmos workflow <name> --file <file>")
	workflowCmd.PersistentFlags().Bool("dry-run", false, "atmos workflow <name> --file <file> --dry-run")
	workflowCmd.PersistentFlags().StringP("stack", "s", "", "atmos workflow <name> --file <file> --stack <stack>")
	workflowCmd.PersistentFlags().String("from-step", "", "atmos workflow <name> --file <file> --from-step <step-name>")

	RootCmd.AddCommand(workflowCmd)
}
