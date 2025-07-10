package zenmoney

import (
	"reflect"
	"testing"
)

func TestResponse_GetIndexedTags(t *testing.T) {
	testCases := []struct {
		name     string
		response *Response
		expected map[string]Tag
	}{
		{
			name: "Single tag",
			response: &Response{
				Tag: []Tag{
					{Id: "tag1", Title: "Tag 1"},
				},
			},
			expected: map[string]Tag{
				"tag1": {Id: "tag1", Title: "Tag 1"},
			},
		},
		{
			name: "Multiple tags",
			response: &Response{
				Tag: []Tag{
					{Id: "tag1", Title: "Tag 1"},
					{Id: "tag2", Title: "Tag 2"},
					{Id: "tag3", Title: "Tag 3"},
				},
			},
			expected: map[string]Tag{
				"tag1": {Id: "tag1", Title: "Tag 1"},
				"tag2": {Id: "tag2", Title: "Tag 2"},
				"tag3": {Id: "tag3", Title: "Tag 3"},
			},
		},
		{
			name: "No tags",
			response: &Response{
				Tag: []Tag{},
			},
			expected: map[string]Tag{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.response.GetIndexedTags()
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("GetIndexedTags() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

func TestResponse_GetIndexedAccounts(t *testing.T) {
	testCases := []struct {
		name     string
		response *Response
		expected map[string]Account
	}{
		{
			name: "Single account",
			response: &Response{
				Account: []Account{
					{Id: "account1", Title: "Account 1"},
				},
			},
			expected: map[string]Account{
				"account1": {Id: "account1", Title: "Account 1"},
			},
		},
		{
			name: "Multiple accounts",
			response: &Response{
				Account: []Account{
					{Id: "account1", Title: "Account 1"},
					{Id: "account2", Title: "Account 2"},
					{Id: "account3", Title: "Account 3"},
				},
			},
			expected: map[string]Account{
				"account1": {Id: "account1", Title: "Account 1"},
				"account2": {Id: "account2", Title: "Account 2"},
				"account3": {Id: "account3", Title: "Account 3"},
			},
		},
		{
			name: "No accounts",
			response: &Response{
				Account: []Account{},
			},
			expected: map[string]Account{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.response.GetIndexedAccounts()
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("GetIndexedAccounts() = %v, expected %v", result, tc.expected)
			}
		})
	}
}
