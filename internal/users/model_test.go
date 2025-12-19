package users

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tCases := []struct {
		name     string
		email    string
		password string
		wantErr  bool
		err      error
	}{
		{
			name:     "good_User",
			email:    "goodEmail@gmai.com",
			password: "qwerty1234",
		},
		{
			name:     "empty_email",
			email:    "",
			password: "qwerty1234",
			wantErr:  true,
			err:      errValidate,
		},
		{
			name:     "empty_pass",
			email:    "goodEmail@gmai.com",
			password: "",
			wantErr:  true,
			err:      errValidate,
		},
		{
			name:     "deleted_email",
			password: "qwerty1234",
			wantErr:  true,
			err:      errValidate,
		},
		{
			name:    "deleted_pass",
			email:   "goodEmail@gmai.com",
			wantErr: true,
			err:     errValidate,
		},
		{
			name:    "0 values",
			wantErr: true,
			err:     errValidate,
		},
	}

	for _, tt := range tCases {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				u, err := NewUser(tt.email, tt.password)
				require.NoError(t, err)
				require.NotNil(t, u)
				require.Equal(t, u.Email, tt.email)
			} else {
				_, err := NewUser(tt.email, tt.password)
				require.EqualError(t, err, tt.err.Error())
			}

		})
	}

}
