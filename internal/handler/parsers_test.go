package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseChatID(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue uint
		expectedError bool
	}{
		{
			name:          "Empty ID",
			input:         "",
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "Negative ID",
			input:         "-1",
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "NaN ID",
			input:         "abc",
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "Float ID",
			input:         "3.14",
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "Valid ID",
			input:         "1",
			expectedValue: 1,
			expectedError: false,
		},
		{
			name:          "Zero ID",
			input:         "0",
			expectedValue: 0,
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			value, err := parseChatID(tc.input)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedValue, value)
			}
		})
	}
}

func TestParseLimit(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue int
	}{
		{
			name:          "Empty limit",
			input:         "",
			expectedValue: limitDefault,
		},
		{
			name:          "Negative limit",
			input:         "-1",
			expectedValue: limitDefault,
		},
		{
			name:          "Too big limit",
			input:         "100",
			expectedValue: limitMax,
		},
		{
			name:          "Zero limit",
			input:         "0",
			expectedValue: limitDefault,
		},
		{
			name:          "NaN limit",
			input:         "abc",
			expectedValue: limitDefault,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			value := parseLimit(tc.input)

			require.Equal(t, tc.expectedValue, value)
		})
	}
}
