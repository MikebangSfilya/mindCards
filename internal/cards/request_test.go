package cards

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		wantError   bool
		expectedErr error
	}{
		{
			name:      "valid_description",
			input:     "this is correct test description, all OKEY",
			wantError: false,
		},
		{
			name:        "empty_NewDescription",
			input:       "",
			wantError:   true,
			expectedErr: errEmptyDescription,
		},
		{
			name:        "only_spaces",
			input:       "     ",
			wantError:   true,
			expectedErr: errEmptyDescription,
		},
		{
			name:        "too_short_9_chars",
			input:       "123456789",
			wantError:   true,
			expectedErr: errShortDescription,
		},
		{
			name:      "exact_min_length_10_chars",
			input:     "1234567890",
			wantError: false,
		},
		{
			name:      "with_whitespace",
			input:     "   valid answer   ",
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			updDto := Update{NewDescription: tc.input}
			err := updDto.Validate()
			if !tc.wantError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
			}

		})
	}

}
