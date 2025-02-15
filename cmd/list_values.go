package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/pkg/config"
	l "github.com/cloudposse/atmos/pkg/list"
	"github.com/cloudposse/atmos/pkg/logger"
	"github.com/cloudposse/atmos/pkg/schema"
	"github.com/cloudposse/atmos/pkg/ui/theme"
	u "github.com/cloudposse/atmos/pkg/utils"
)

// listValuesCmd lists component values across stacks
var listValuesCmd = &cobra.Command{
	Use:   "values [component]",
	Short: "List component values across stacks",
	Long:  "List values for a component across all stacks where it is used",
	Example: "atmos list values vpc\n" +
		"atmos list values vpc --query .vars\n" +
		"atmos list values vpc --abstract\n" +
		"atmos list values vpc --max-columns 5\n" +
		"atmos list values vpc --format json\n" +
		"atmos list values vpc --format yaml\n" +
		"atmos list values vpc --format csv",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check Atmos configuration
		checkAtmosConfig()

		// Initialize logger from CLI config
		configAndStacksInfo := schema.ConfigAndStacksInfo{}
		atmosConfig, err := config.InitCliConfig(configAndStacksInfo, true)
		if err != nil {
			fmt.Printf("Error initializing CLI config: %v\n", err)
			return
		}

		log, err := logger.NewLoggerFromCliConfig(atmosConfig)
		if err != nil {
			fmt.Printf("Error initializing logger: %v\n", err)
			return
		}

		flags := cmd.Flags()

		queryFlag, err := flags.GetString("query")
		if err != nil {
			log.Error(fmt.Errorf("failed to get query flag: %w", err))
			return
		}

		abstractFlag, err := flags.GetBool("abstract")
		if err != nil {
			log.Error(fmt.Errorf("failed to get abstract flag: %w", err))
			return
		}

		maxColumnsFlag, err := flags.GetInt("max-columns")
		if err != nil {
			log.Error(fmt.Errorf("failed to get max-columns flag: %w", err))
			return
		}

		formatFlag, err := flags.GetString("format")
		if err != nil {
			log.Error(fmt.Errorf("failed to get format flag: %w", err))
			return
		}

		delimiterFlag, err := flags.GetString("delimiter")
		if err != nil {
			log.Error(fmt.Errorf("failed to get delimiter flag: %w", err))
			return
		}

		// Set appropriate default delimiter based on format
		if formatFlag == l.FormatCSV && delimiterFlag == l.DefaultTSVDelimiter {
			delimiterFlag = l.DefaultCSVDelimiter
		}

		component := args[0]

		// Get all stacks
		stacksMap, err := e.ExecuteDescribeStacks(atmosConfig, "", nil, nil, nil, false, false, false, false, nil)
		if err != nil {
			log.Error(fmt.Errorf("failed to describe stacks: %w", err))
			return
		}

		output, err := l.FilterAndListValues(stacksMap, component, queryFlag, abstractFlag, maxColumnsFlag, formatFlag, delimiterFlag)
		if err != nil {
			log.Warning(fmt.Sprintf("Failed to filter and list values: %v", err))
			return
		}

		log.Info(output)
	},
}

// listVarsCmd is an alias for 'list values --query .vars'
var listVarsCmd = &cobra.Command{
	Use:   "vars [component]",
	Short: "List component vars across stacks (alias for 'list values --query .vars')",
	Long:  "List vars for a component across all stacks where it is used",
	Example: "atmos list vars vpc\n" +
		"atmos list vars vpc --abstract\n" +
		"atmos list vars vpc --max-columns 5\n" +
		"atmos list vars vpc --format json\n" +
		"atmos list vars vpc --format yaml\n" +
		"atmos list vars vpc --format csv",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Set the query flag to .vars
		if err := cmd.Flags().Set("query", ".vars"); err != nil {
			u.PrintMessageInColor(fmt.Sprintf("Error setting query flag: %v", err), theme.Colors.Error)
			return
		}
		// Run the values command
		listValuesCmd.Run(cmd, args)
	},
}

func init() {
	// Flags for both commands
	commonFlags := func(cmd *cobra.Command) {
		cmd.PersistentFlags().String("query", "", "JMESPath query to filter values")
		cmd.PersistentFlags().Bool("abstract", false, "Include abstract components")
		cmd.PersistentFlags().Int("max-columns", 10, "Maximum number of columns to display")
		cmd.PersistentFlags().String("format", "", "Output format (table, json, yaml, csv, tsv)")
		cmd.PersistentFlags().String("delimiter", "\t", "Delimiter for csv/tsv output (default: tab for tsv, comma for csv)")
	}

	commonFlags(listValuesCmd)
	commonFlags(listVarsCmd)

	listCmd.AddCommand(listValuesCmd)
	listCmd.AddCommand(listVarsCmd)
}
