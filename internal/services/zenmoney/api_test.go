package zenmoney

import (
	"errors"
	"net/http"
	"os"
	"testing"
)

func TestNewApi(t *testing.T) {
	client := &http.Client{}
	api := NewApi(client)

	if api == nil {
		t.Error("NewApi() returned nil")
	}

}

func TestApi_Init(t *testing.T) {
	client := &http.Client{}
	api := &Api{client: client}

	tests := []struct {
		name          string
		envToken      string
		expectedToken string
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "Valid token",
			envToken:      "test-token-123",
			expectedToken: "test-token-123",
			expectedError: false,
		},
		{
			name:          "Empty token",
			envToken:      "",
			expectedToken: "",
			expectedError: true,
			errorMessage:  "you need to set ZENMONEY TOKEN environment variable",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			originalToken := os.Getenv("ZENMONEY_TOKEN")

			if tc.envToken != "" {
				os.Setenv("ZENMONEY_TOKEN", tc.envToken)
			} else {
				os.Unsetenv("ZENMONEY_TOKEN")
			}

			defer func() {
				if originalToken != "" {
					os.Setenv("ZENMONEY_TOKEN", originalToken)
				} else {
					os.Unsetenv("ZENMONEY_TOKEN")
				}
			}()

			token, err := api.Init()

			if tc.expectedError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				if err.Error() != tc.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tc.errorMessage, err.Error())
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token '%s', got '%s'", tc.expectedToken, token)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token != tc.expectedToken {
					t.Errorf("Expected token '%s', got '%s'", tc.expectedToken, token)
				}
			}
		})
	}
}

func TestApi_Diff_NotImplemented(t *testing.T) {
	client := &http.Client{}
	api := &Api{client: client}

	var _ ApiInterface = api
}

func TestConstants(t *testing.T) {
	expectedBaseURL := "https://api.zenmoney.app/v8/diff/"
	if BASE_URL != expectedBaseURL {
		t.Errorf("Expected BASE_URL to be '%s', got '%s'", expectedBaseURL, BASE_URL)
	}
}

type MockApi struct {
	DiffFunc func() (*Response, error)
}

func (m *MockApi) Diff() (*Response, error) {
	if m.DiffFunc != nil {
		return m.DiffFunc()
	}
	return nil, errors.New("not implemented")
}

func TestApiInterface(t *testing.T) {

	mockResponse := &Response{}
	mockApi := &MockApi{
		DiffFunc: func() (*Response, error) {
			return mockResponse, nil
		},
	}

	result, err := mockApi.Diff()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != mockResponse {
		t.Error("Expected mock response")
	}

	mockError := errors.New("api error")
	mockApi = &MockApi{
		DiffFunc: func() (*Response, error) {
			return nil, mockError
		},
	}

	result, err = mockApi.Diff()
	if err == nil {
		t.Error("Expected error but got nil")
	}
	if err.Error() != mockError.Error() {
		t.Errorf("Expected error '%s', got '%s'", mockError.Error(), err.Error())
	}
	if result != nil {
		t.Error("Expected nil result")
	}
}
