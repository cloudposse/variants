package config

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/cloudposse/atmos/pkg/schema"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Successfully unmarshal valid config data from viper into atmosConfig struct
func TestDeepMergeConfigUnmarshalValidConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigType("yaml")

	validConfig := []byte(`
    stacks:
        base_path: "stacks"
    `)

	err := v.ReadConfig(bytes.NewBuffer(validConfig))
	require.NoError(t, err)

	cl := &ConfigLoader{
		viper:       v,
		atmosConfig: schema.AtmosConfiguration{},
	}

	err = cl.deepMergeConfig()
	require.NoError(t, err)

	require.Equal(t, "stacks", cl.atmosConfig.Stacks.BasePath)
}

// Returns list of atmos config files with supported extensions  .yaml, .yml
func TestSearchAtmosConfigFileDir_ReturnsConfigFilesWithSupportedExtensions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "atmos-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with different extensions
	files := []string{
		"atmos.yaml",
		"atmos.yml",
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	viper := viper.New()
	viper.SetConfigType("yaml")
	cl := &ConfigLoader{
		viper:       viper,
		atmosConfig: schema.AtmosConfiguration{},
	}
	got, err := cl.SearchAtmosConfig(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 config files, got %d", len(got))
	}

	// Verify extensions are in correct order
	if !strings.HasSuffix(got[0], "atmos.yaml") {
		t.Errorf("Expected config files with supported extensions, got %v", got)
	}
}

// Successfully load single config file from valid command line argument
func TestLoadExplicitConfigsWithValidConfigFile(t *testing.T) {
	// Setup test config file
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	configPath := filepath.Join(tmpDir, "atmos.yaml")

	err := os.WriteFile(configPath, []byte("test: config"), 0o644)
	require.NoError(t, err)

	cl := &ConfigLoader{
		atmosConfig: schema.AtmosConfiguration{},
		viper:       viper.New(),
	}

	err = cl.loadExplicitConfigs([]string{configPath})
	require.NoError(t, err)

	assert.Contains(t, cl.AtmosConfigPaths, configPath)
}

// Successfully load multiple config file from valid command line argument and directories
func TestLoadExplicitConfigsWithMultipleConfigFiles(t *testing.T) {
	// Setup test config files
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	configPath1 := filepath.Join(tmpDir, "atmos.yaml")
	configPath2 := filepath.Join(tmpDir, "atmos.yml")
	err := os.WriteFile(configPath1, []byte("test: config1"), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(configPath2, []byte("test: config2"), 0o644)
	require.NoError(t, err)
	cl := &ConfigLoader{
		atmosConfig: schema.AtmosConfiguration{},
		viper:       viper.New(),
	}

	err = cl.loadExplicitConfigs([]string{configPath1, configPath2})
	require.NoError(t, err)
	assert.Contains(t, cl.AtmosConfigPaths, configPath1)
	assert.Contains(t, cl.AtmosConfigPaths, configPath2)
	// test read from dir
	cl = &ConfigLoader{
		atmosConfig: schema.AtmosConfiguration{},
		viper:       viper.New(),
	}
	err = cl.loadExplicitConfigs([]string{tmpDir})
	require.NoError(t, err)
	assert.Contains(t, cl.AtmosConfigPaths, tmpDir)
}

// Function correctly prioritizes .yaml over .yml for same base filename
func TestDetectPriorityFilesPreferYamlOverYml(t *testing.T) {
	cl := &ConfigLoader{}

	files := []string{
		"config/app.yml",
		"config/app.yaml",
		"config/db.yml",
	}

	result := cl.detectPriorityFiles(files)

	expected := []string{
		"config/app.yaml",
		"config/db.yml",
	}

	result = cl.sortFilesByDepth(result)
	expected = cl.sortFilesByDepth(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

// Sort files by directory depth in ascending order
func TestSortFilesByDepthSortsFilesCorrectly(t *testing.T) {
	cl := &ConfigLoader{}

	files := []string{
		"a/b/c/file1.yaml",
		"x/file2.yaml",
		"file1.yaml",
		"a/b/file3.yaml",
		"file4.yaml",
	}

	expected := []string{
		"file1.yaml",
		"file4.yaml",
		"x/file2.yaml",
		"a/b/file3.yaml",
		"a/b/c/file1.yaml",
	}

	result := cl.sortFilesByDepth(files)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestDownloadRemoteConfig(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("mock content"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer mockServer.Close()
	viper := viper.New()
	viper.SetConfigType("yaml")
	cl := &ConfigLoader{
		viper:       viper,
		atmosConfig: schema.AtmosConfiguration{},
	}
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	t.Run("Valid URL", func(t *testing.T) {
		tempFile, err := cl.downloadRemoteConfig(mockServer.URL, tmpDir)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify the temporary file contains the correct content
		content, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}
		if string(content) != "mock content" {
			t.Errorf("Unexpected file content: got %v, want %v", string(content), "mock content")
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		_, err := cl.downloadRemoteConfig("http://invalid-url", tmpDir)
		if err == nil {
			t.Fatal("Expected an error for invalid URL, got nil")
		}
	})
}

func TestParseArraySeparatorWithMultipleParts(t *testing.T) {
	input := "part1; part2;part3 ;part4"

	result := parseArraySeparator(input)

	expected := []string{"part1", "part2", "part3", "part4"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d but got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %s at position %d but got %s", v, i, result[i])
		}
	}
	input = "part1"
	result = parseArraySeparator(input)
	expected = []string{"part1"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d but got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %s at position %d but got %s", v, i, result[i])
		}
	}
}
