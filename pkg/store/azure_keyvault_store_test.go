package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKeyVaultClient is a mock implementation of the Azure Key Vault client
type MockKeyVaultClient struct {
	mock.Mock
}

func (m *MockKeyVaultClient) SetSecret(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (azsecrets.SetSecretResponse, error) {
	args := m.Called(ctx, name, parameters)
	return azsecrets.SetSecretResponse{}, args.Error(1)
}

func (m *MockKeyVaultClient) GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return azsecrets.GetSecretResponse{}, args.Error(1)
	}
	value := args.String(0)
	return azsecrets.GetSecretResponse{Value: &value}, args.Error(1)
}

func TestNewKeyVaultStore(t *testing.T) {
	tests := []struct {
		name      string
		options   KeyVaultStoreOptions
		wantError bool
	}{
		{
			name: "valid options",
			options: KeyVaultStoreOptions{
				VaultURL: "https://test-vault.vault.azure.net/",
			},
			wantError: false,
		},
		{
			name: "missing vault url",
			options: KeyVaultStoreOptions{
				VaultURL: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKeyVaultStore(tt.options)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestKeyVaultStore_getKey(t *testing.T) {
	delimiter := "-"
	store := &KeyVaultStore{
		prefix:         "prefix",
		stackDelimiter: &delimiter,
	}

	tests := []struct {
		name      string
		stack     string
		component string
		key       string
		expected  string
	}{
		{
			name:      "simple path",
			stack:     "dev",
			component: "app",
			key:       "config",
			expected:  "prefix-dev-app-config",
		},
		{
			name:      "nested component",
			stack:     "dev",
			component: "app/service",
			key:       "config",
			expected:  "prefix-dev-app-service-config",
		},
		{
			name:      "multi-level stack",
			stack:     "dev-us-west-2",
			component: "app",
			key:       "config",
			expected:  "prefix-dev-us-west-2-app-config",
		},
		{
			name:      "uppercase characters",
			stack:     "Dev",
			component: "App/Service",
			key:       "Config",
			expected:  "prefix-dev-app-service-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.getKey(tt.stack, tt.component, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKeyVaultStore_InputValidation(t *testing.T) {
	mockClient := new(MockKeyVaultClient)
	delimiter := "-"
	store := &KeyVaultStore{
		client:         mockClient,
		prefix:         "prefix",
		stackDelimiter: &delimiter,
	}

	tests := []struct {
		name      string
		stack     string
		component string
		key       string
		value     interface{}
		operation string
		mockFn    func()
		wantError bool
	}{
		{
			name:      "empty stack",
			stack:     "",
			component: "app",
			key:       "config",
			value:     "test",
			operation: "set",
			mockFn:    func() {},
			wantError: true,
		},
		{
			name:      "empty component",
			stack:     "dev",
			component: "",
			key:       "config",
			value:     "test",
			operation: "set",
			mockFn:    func() {},
			wantError: true,
		},
		{
			name:      "empty key",
			stack:     "dev",
			component: "app",
			key:       "",
			value:     "test",
			operation: "set",
			mockFn:    func() {},
			wantError: true,
		},
		{
			name:      "non-string value",
			stack:     "dev",
			component: "app",
			key:       "config",
			value:     123,
			operation: "set",
			mockFn:    func() {},
			wantError: true,
		},
		{
			name:      "valid set operation",
			stack:     "dev",
			component: "app",
			key:       "config",
			value:     "test",
			operation: "set",
			mockFn: func() {
				mockClient.On("SetSecret", mock.Anything, "prefix-dev-app-config", mock.Anything).
					Return(azsecrets.SetSecretResponse{}, nil)
			},
			wantError: false,
		},
		{
			name:      "valid get operation",
			stack:     "dev",
			component: "app",
			key:       "config",
			operation: "get",
			mockFn: func() {
				mockClient.On("GetSecret", mock.Anything, "prefix-dev-app-config").
					Return("test-value", nil)
			},
			wantError: false,
		},
		{
			name:      "get operation error",
			stack:     "dev",
			component: "app",
			key:       "config",
			operation: "get",
			mockFn: func() {
				mockClient.On("GetSecret", mock.Anything, "prefix-dev-app-config").
					Return(nil, fmt.Errorf("secret not found"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockClient.Calls = nil
			tt.mockFn()

			var err error
			if tt.operation == "set" {
				err = store.Set(tt.stack, tt.component, tt.key, tt.value)
			} else {
				_, err = store.Get(tt.stack, tt.component, tt.key)
			}

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
