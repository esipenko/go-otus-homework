package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	StrLenTest struct {
		Str string `validate:"len:10"`
	}

	StrInTest struct {
		Str string `validate:"in:a,b"`
	}

	StrRegexpTest struct {
		Str string `validate:"regexp:^[a-zA-Z]+$"`
	}

	StrComplexTest struct {
		Str string `validate:"regexp:^[a-zA-Z]+$|len:3|in:aaa"`
	}

	StrUnsupportedTest struct {
		Str string `validate:"min:12"`
	}

	StrSliceComplexValidation struct {
		Str []string `validate:"regexp:^[a-zA-Z]+$|len:3|in:aaa,bbb"`
	}

	IntMinTest struct {
		Int int `validate:"min:12"`
	}

	IntMaxTest struct {
		Int int `validate:"max:12"`
	}

	IntInTest struct {
		Int int `validate:"in:12,13,14"`
	}

	IntComplexTest struct {
		Int int `validate:"in:12,13,14|min:12|max:14"`
	}

	IntUnsupportedTest struct {
		Int int `validate:"regexp:^[a-zA-Z]+$"`
	}

	IntSliceTest struct {
		Int []int `validate:"in:12,13,14|min:12|max:14"`
	}

	User struct {
		ID     string `json:"id" validate:"len:10"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte `validate:"len:10"`
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestStrRules(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		len         int
	}{
		{
			StrLenTest{Str: "1"},
			ErrLenError,
			1,
		},
		{
			StrLenTest{Str: "0123456789"},
			nil,
			0,
		},
		{
			StrInTest{Str: "c"},
			ErrInError,
			1,
		},
		{
			StrInTest{Str: "a"},
			nil,
			0,
		},
		{
			StrRegexpTest{Str: "1"},
			ErrRegexpError,
			1,
		},
		{
			StrRegexpTest{Str: "a"},
			nil,
			0,
		},
		{
			StrComplexTest{Str: "1234"},
			ValidationErrors{},
			3,
		},
		{
			StrComplexTest{Str: "aaa"},
			nil,
			0,
		},
		{
			StrUnsupportedTest{Str: "aaa"},
			ErrUnsupportedValidation,
			1,
		},
		{
			StrSliceComplexValidation{Str: []string{"aaa", "1234"}},
			ValidationErrors{},
			3,
		},
		{
			StrSliceComplexValidation{Str: []string{"aaa", "bbb"}},
			nil,
			0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			fmt.Println(err)

			if tt.expectedErr != nil {
				var validationErrors ValidationErrors
				if errors.As(err, &validationErrors) {
					require.Equal(t, tt.len, validationErrors.len())
				}
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			_ = tt
		})
	}
}

func TestIntRules(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		len         int
	}{
		{
			IntMinTest{Int: 0},
			ErrMinError,
			1,
		},
		{
			IntMinTest{Int: 1000},
			nil,
			0,
		},
		{
			IntMaxTest{Int: 100},
			ErrMaxError,
			1,
		},
		{
			IntMaxTest{Int: 0},
			nil,
			0,
		},
		{
			IntInTest{Int: 0},
			ErrInError,
			1,
		},
		{
			IntInTest{Int: 12},
			nil,
			0,
		},
		{
			IntComplexTest{Int: -1},
			ValidationErrors{},
			2,
		},
		{
			IntComplexTest{Int: 12},
			nil,
			0,
		},
		{
			IntUnsupportedTest{Int: 0},
			ErrUnsupportedValidation,
			1,
		},
		{
			IntSliceTest{Int: []int{-1, -1234}},
			ValidationErrors{},
			4,
		},
		{
			IntSliceTest{Int: []int{12, 13}},
			nil,
			0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			fmt.Println(err)

			if tt.expectedErr != nil {
				var validationErrors ValidationErrors
				if errors.As(err, &validationErrors) {
					require.Equal(t, tt.len, validationErrors.len())
				}
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			_ = tt
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		len         int
	}{
		{
			make([]any, 0),
			ErrUnsupportedType,
			0,
		},
		{
			Token{Header: make([]byte, 0)},
			ErrUnsupportedType,
			0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			err := Validate(tt.in)
			fmt.Println(err)

			if tt.expectedErr != nil {
				var validationErrors ValidationErrors
				if errors.As(err, &validationErrors) {
					require.Equal(t, tt.len, validationErrors.len())
				}
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			_ = tt
		})
	}
}
