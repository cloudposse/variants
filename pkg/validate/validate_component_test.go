package validate

import (
	e "github.com/cloudposse/atmos/internal/exec"
	c "github.com/cloudposse/atmos/pkg/config"
	u "github.com/cloudposse/atmos/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateComponent(t *testing.T) {
	Config, err := c.InitConfig(c.ConfigAndStacksInfo{})
	assert.Nil(t, err)

	_, err = e.ExecuteValidateComponent(Config, "infra/vpc", "tenant1-ue2-dev", "validate-infra-vpc-component.rego", "opa")
	u.PrintError(err)
	assert.Error(t, err)
}

func TestValidateComponent2(t *testing.T) {
	Config, err := c.InitConfig(c.ConfigAndStacksInfo{})
	assert.Nil(t, err)

	_, err = e.ExecuteValidateComponent(Config, "infra/vpc", "tenant1-ue2-prod", "", "")
	u.PrintError(err)
	assert.Error(t, err)
}

func TestValidateComponent3(t *testing.T) {
	Config, err := c.InitConfig(c.ConfigAndStacksInfo{})
	assert.Nil(t, err)

	_, err = e.ExecuteValidateComponent(Config, "infra/vpc", "tenant1-ue2-staging", "", "")
	u.PrintError(err)
	assert.Error(t, err)
}
