package hw09structvalidator

import (
	"reflect"
)

type sliceForValidate struct {
	Field string
	Value reflect.Value
	Kind  reflect.Kind
	Rules validationRules
}

func validateSlice(slice sliceForValidate) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)

	//exhaustive:ignore
	switch slice.Kind {
	case reflect.Int:
		for i := 0; i < slice.Value.Len(); i++ {
			value := slice.Value.Index(i)
			errors, err := validateIntRules(slice.Field, int(value.Int()), slice.Rules)
			if err != nil {
				return nil, err
			}
			validationErrors = append(validationErrors, errors...)
		}
	case reflect.String:
		for i := 0; i < slice.Value.Len(); i++ {
			value := slice.Value.Index(i)
			errors, err := validateStringRules(slice.Field, value.String(), slice.Rules)
			if err != nil {
				return nil, err
			}
			validationErrors = append(validationErrors, errors...)
		}
	}

	return validationErrors, nil
}
