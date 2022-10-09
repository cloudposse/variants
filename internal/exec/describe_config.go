package exec

import (
	c "github.com/cloudposse/atmos/pkg/config"
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ExecuteDescribeConfig executes `describe config` command
func ExecuteDescribeConfig(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	format, err := flags.GetString("format")
	if err != nil {
		return err
	}

	cliConfig, err := c.InitCliConfig(c.ConfigAndStacksInfo{}, true)
	if err != nil {
		return err
	}

	if format == "json" {
		err = u.PrintAsJSON(cliConfig)
	} else if format == "yaml" {
		err = u.PrintAsYAML(cliConfig)
	} else {
		err = errors.New("invalid flag '--format'. Accepted values are 'json' or 'yaml'")
	}
	if err != nil {
		return err
	}

	return nil
}
