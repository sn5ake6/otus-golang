package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrInvalidType = errors.New("invalid type")

type ValidationError struct {
	Field string
	Err   error
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("Field: %v, Error: %v", err.Field, err.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errors := make([]string, 0)
	for _, err := range v {
		errors = append(errors, err.Error())
	}

	return strings.Join(errors, "\n")
}

func Validate(v interface{}) error {
	iv := reflect.ValueOf(v)

	if iv.Kind() != reflect.Struct {
		return ErrInvalidType
	}

	it := iv.Type()

	validationErrors := make(ValidationErrors, 0)

	for i := 0; i < it.NumField(); i++ {
		field := it.Field(i)

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		rules := parseRules(validateTag)

		fieldValue := iv.Field(i)

		var errors ValidationErrors
		var err error

		//exhaustive:ignore
		switch field.Type.Kind() {
		case reflect.Int:
			errors, err = validateIntRules(field.Name, int(fieldValue.Int()), rules)
			if err != nil {
				return err
			}
		case reflect.String:
			errors, err = validateStringRules(field.Name, fieldValue.String(), rules)
			if err != nil {
				return err
			}
		case reflect.Slice:
			slice := sliceForValidate{
				Field: field.Name,
				Value: fieldValue,
				Kind:  field.Type.Elem().Kind(),
				Rules: rules,
			}
			errors, err = validateSlice(slice)
			if err != nil {
				return err
			}
		}
		validationErrors = append(validationErrors, errors...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
