package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/elewis787/boa"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"

	e "github.com/cloudposse/atmos/internal/exec"
	"github.com/cloudposse/atmos/internal/tui/templates"
	tuiUtils "github.com/cloudposse/atmos/internal/tui/utils"
	cfg "github.com/cloudposse/atmos/pkg/config"
	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

var atmosConfig schema.AtmosConfiguration

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "atmos",
	Short: "Universal Tool for DevOps and Cloud Automation",
	Long:  `Atmos is a universal tool for DevOps and cloud automation used for provisioning, managing and orchestrating workflows across various toolchains`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Determine if the command is a help command or if the help flag is set
		isHelpCommand := cmd.Name() == "help"
		helpFlag := cmd.Flags().Changed("help")

		isHelpRequested := isHelpCommand || helpFlag

		if isHelpRequested {
			// Do not silence usage or errors when help is invoked
			cmd.SilenceUsage = false
			cmd.SilenceErrors = false
		} else {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check Atmos configuration
		checkAtmosConfig()

		// Print a styled Atmos logo to the terminal
		fmt.Println()
		err := tuiUtils.PrintStyledText("ATMOS")
		if err != nil {
			u.LogErrorAndExit(err)
		}

		err = e.ExecuteAtmosCmd()
		if err != nil {
			u.LogErrorAndExit(err)
		}
	},
}

func setupLogger(atmosConfig *schema.AtmosConfiguration) {
	switch atmosConfig.Logs.Level {
	case "Trace":
		log.SetLevel(log.DebugLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warning":
		log.SetLevel(log.WarnLevel)
	case "Off":
		log.SetOutput(io.Discard)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.Info("log level set from atmos config", "level", log.GetLevel())

	if atmosConfig.Logs.File != "/dev/stderr" {
		log.Info("setting log file", "file", atmosConfig.Logs.File)
		logFile, err := os.OpenFile(atmosConfig.Logs.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiGreen + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	// InitCliConfig finds and merges CLI configurations in the following order:
	// system dir, home dir, current dir, ENV vars, command-line arguments
	// Here we need the custom commands from the config
	var initErr error
	atmosConfig, initErr = cfg.InitCliConfig(schema.ConfigAndStacksInfo{}, false)
	if initErr != nil && !errors.Is(initErr, cfg.NotFound) {
		if isVersionCommand() {
			log.Debug("warning: CLI configuration 'atmos.yaml' file not found", "error", initErr)
		} else {
			u.LogErrorAndExit(initErr)
		}
	}

	// Set the log level for the charmbracelet/log package based on the atmosConfig
	setupLogger(&atmosConfig)

	var err error
	// If CLI configuration was found, process its custom commands and command aliases
	if initErr == nil {
		err = processCustomCommands(atmosConfig, atmosConfig.Commands, RootCmd, true)
		if err != nil {
			u.LogErrorAndExit(err)
		}

		err = processCommandAliases(atmosConfig, atmosConfig.CommandAliases, RootCmd, true)
		if err != nil {
			u.LogErrorAndExit(err)
		}
	}

	return RootCmd.Execute()
}

func init() {
	// Add template function for wrapped flag usages
	cobra.AddTemplateFunc("wrappedFlagUsages", templates.WrappedFlagUsages)

	RootCmd.PersistentFlags().String("redirect-stderr", "", "File descriptor to redirect 'stderr' to. "+
		"Errors can be redirected to any file or any standard file descriptor (including '/dev/null'): atmos <command> --redirect-stderr /dev/stdout")

	RootCmd.PersistentFlags().String("logs-level", "Info", "Logs level. Supported log levels are Trace, Debug, Info, Warning, Off. If the log level is set to Off, Atmos will not log any messages")
	RootCmd.PersistentFlags().String("logs-file", "/dev/stdout", "The file to write Atmos logs to. Logs can be written to any file or any standard file descriptor, including '/dev/stdout', '/dev/stderr' and '/dev/null'")

	// Set custom usage template
	err := templates.SetCustomUsageFunc(RootCmd)
	if err != nil {
		u.LogErrorAndExit(err)
	}

	initCobraConfig()
}

func initCobraConfig() {
	RootCmd.SetOut(os.Stdout)
	styles := boa.DefaultStyles()
	b := boa.New(boa.WithStyles(styles))
	oldUsageFunc := RootCmd.UsageFunc()
	RootCmd.SetUsageFunc(b.UsageFunc)

	RootCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		// Print a styled Atmos logo to the terminal
		fmt.Println()
		if command.Use != "atmos" || command.Flags().Changed("help") {
			err := tuiUtils.PrintStyledText("ATMOS")
			if err != nil {
				u.LogErrorAndExit(err)
			}
			if err := oldUsageFunc(command); err != nil {
				u.LogErrorAndExit(err)
			}
		} else {
			err := tuiUtils.PrintStyledText("ATMOS")
			if err != nil {
				u.LogErrorAndExit(err)
			}
			b.HelpFunc(command, strings)
			if err := command.Usage(); err != nil {
				u.LogErrorAndExit(err)
			}
		}
		CheckForAtmosUpdateAndPrintMessage(atmosConfig)
	})
}

// https://www.sobyte.net/post/2021-12/create-cli-app-with-cobra/
// https://github.com/spf13/cobra/blob/master/user_guide.md
// https://blog.knoldus.com/create-kubectl-like-cli-with-go-and-cobra/
// https://pkg.go.dev/github.com/c-bata/go-prompt
// https://pkg.go.dev/github.com/spf13/cobra
// https://scene-si.org/2017/04/20/managing-configuration-with-viper/
