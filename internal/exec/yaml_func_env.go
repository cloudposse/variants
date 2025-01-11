package exec

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudposse/atmos/pkg/config"
	"github.com/cloudposse/atmos/pkg/schema"
	u "github.com/cloudposse/atmos/pkg/utils"
)

func processTagEnv(
	atmosConfig schema.AtmosConfiguration,
	input string,
	currentStack string,
) any {
	u.LogTrace(atmosConfig, fmt.Sprintf("Executing Atmos YAML function: %s", input))

	str, err := getStringAfterTag(input, config.AtmosYamlFuncEnv)

	if err != nil {
		u.LogErrorAndExit(atmosConfig, err)
	}

	res := os.Getenv(str)
	if err != nil {
		u.LogErrorAndExit(atmosConfig, err)
	}

	var decoded any
	if err = json.Unmarshal([]byte(res), &decoded); err != nil {
		return res
	}

	return decoded
}
