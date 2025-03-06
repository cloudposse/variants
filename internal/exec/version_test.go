package exec

import (
	"testing"

	"github.com/cloudposse/atmos/pkg/version"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

// Mock interfaces generated by mockgen
// VersionExecutor defines the interface for version execution operations
//
//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
type VersionExecutor interface {
	PrintStyledText(text string) error
	GetLatestGitHubRepoRelease(owner, repo string) (string, error)
	PrintMessage(message string)
	PrintMessageToUpgradeToAtmosLatestRelease(version string)
}

func TestVersionExec_Execute(t *testing.T) {
	// Save original values
	originalVersion := version.Version
	defer func() { version.Version = originalVersion }()

	tests := []struct {
		name                string
		checkFlag           bool
		version             string
		latestRelease       string
		printStyledTextErr  error
		getLatestReleaseErr error
	}{
		{
			name:      "Basic execution without check",
			checkFlag: false,
			version:   "v1.0.0",
		},
		{
			name:          "Check flag with same version",
			checkFlag:     true,
			version:       "v1.0.0",
			latestRelease: "v1.0.0",
		},
		{
			name:          "Check flag with different version",
			checkFlag:     true,
			version:       "v1.0.0",
			latestRelease: "v1.1.0",
		},
		{
			name:                "Check flag with release check error",
			checkFlag:           true,
			version:             "v1.0.0",
			getLatestReleaseErr: errors.New("github error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
			mockExec := NewMockVersionExecutor(ctrl)

			// Set version
			version.Version = tt.version

			// Setup mock expectations
			mockExec.EXPECT().PrintStyledText("ATMOS").Return(tt.printStyledTextErr)
			if tt.printStyledTextErr == nil {
				mockExec.EXPECT().PrintMessage(gomock.Any()).Times(3)

				if tt.checkFlag {
					mockExec.EXPECT().GetLatestGitHubRepoRelease("cloudposse", "atmos").
						Return(tt.latestRelease, tt.getLatestReleaseErr)

					if tt.getLatestReleaseErr == nil && tt.latestRelease != "" {
						if tt.latestRelease != tt.version {
							mockExec.EXPECT().PrintMessageToUpgradeToAtmosLatestRelease(
								gomock.Eq(tt.latestRelease[1:])).Times(1)
						}
					}
				}
			}

			// Create test instance with mocks
			v := versionExec{
				printStyledText:                           mockExec.PrintStyledText,
				getLatestGitHubRepoRelease:                mockExec.GetLatestGitHubRepoRelease,
				printMessage:                              mockExec.PrintMessage,
				printMessageToUpgradeToAtmosLatestRelease: mockExec.PrintMessageToUpgradeToAtmosLatestRelease,
			}

			// Execute the function
			v.Execute(tt.checkFlag)
		})
	}
}
