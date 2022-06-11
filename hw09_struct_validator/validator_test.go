package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Redirect struct {
		Code []int `validate:"in:301,302,307"`
	}

	InvalidTagValueForInt struct {
		Count int `validate:"min:one"`
	}

	InvalidTagValueForString struct {
		Name string `validate:"len:two"`
	}
)

func TestValidateWithInvalidType(t *testing.T) {
	tests := []interface{}{
		5,
		"string",
		0.00,
		false,
		[]int{155},
		[]string{"string"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)

			assert.Equal(t, ErrInvalidType, err)
		})
	}
}

func TestValidateWithValidationErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: App{Version: "V1"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrStringValueLen},
			},
		},
		{
			in: App{Version: "V123456"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrStringValueLen},
			},
		},
		{
			in: Response{Code: 502},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntValueOutOfList},
			},
		},
		{
			in: User{
				ID:    "1234567890",
				Age:   10,
				Email: "test.com",
				Role:  "user",
				Phones: []string{
					"12345",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStringValueLen},
				ValidationError{Field: "Age", Err: ErrIntValueTooSmall},
				ValidationError{Field: "Email", Err: ErrStringValueDoesNotMatchRegexp},
				ValidationError{Field: "Role", Err: ErrStringValueOutOfList},
				ValidationError{Field: "Phones", Err: ErrStringValueLen},
			},
		},
		{
			in: User{
				ID:    "123456789012354678901234567890123456",
				Age:   51,
				Email: "test@test.com",
				Role:  "stuff",
				Phones: []string{
					"12345678901",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrIntValueTooLarge},
			},
		},
		{
			in: Redirect{
				Code: []int{301, 302, 307, 404, 505},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIntValueOutOfList},
				ValidationError{Field: "Code", Err: ErrIntValueOutOfList},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			assert.True(t, errors.As(err, &tt.expectedErr))
			assert.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

func TestValidateWithProgramErrors(t *testing.T) {
	tests := []interface{}{
		InvalidTagValueForInt{Count: 15},
		InvalidTagValueForString{Name: "Name"},
	}

	var validationError *ValidationError

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)

			assert.False(t, errors.As(err, &validationError))
		})
	}
}

func TestValidateWithoutErrors(t *testing.T) {
	tests := []interface{}{
		App{Version: "V1234"},
		Token{},
		Response{Code: 200},
		Response{Code: 404},
		Response{Code: 500},
		User{
			ID:    "123456789012354678901234567890123456",
			Age:   18,
			Email: "test@test.com",
			Role:  "admin",
			Phones: []string{
				"12345678901",
				"55555678901",
			},
		},
		User{
			ID:    "123456789012354678901234567890123456",
			Age:   50,
			Email: "test@test.com",
			Role:  "stuff",
			Phones: []string{
				"12345678901",
			},
		},
		Redirect{
			Code: []int{301, 302, 307},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)

			assert.Nil(t, err)
		})
	}
}
