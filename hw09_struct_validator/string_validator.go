package hw09structvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrStringValueLen = errors.New("invalid string len")

var ErrStringValueDoesNotMatchRegexp = errors.New("value does not match expession")

var ErrStringValueOutOfList = errors.New("value out of list")

func validateStringRules(field string, value string, rules validationRules) (ValidationErrors, error) {
	errors := make(ValidationErrors, 0)
	for _, rule := range rules {
		validateErr, err := validateString(value, rule)
		if err != nil {
			return nil, err
		}

		if validateErr != nil {
			errors = append(errors, ValidationError{Field: field, Err: validateErr})
		}
	}

	return errors, nil
}

func validateString(value string, rule validationRule) (error, error) {
	switch rule.Name {
	case "len":
		expectedLen, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}

		if len(value) != expectedLen {
			return ErrStringValueLen, nil
		}
	case "regexp":
		compiledRegexp, err := regexp.Compile(rule.Value)
		if err != nil {
			return nil, err
		}

		if !compiledRegexp.Match([]byte(value)) {
			return ErrStringValueDoesNotMatchRegexp, nil
		}
	case "in":
		values := strings.Split(rule.Value, ",")
		isIn := false
		for _, v := range values {
			if value == v {
				isIn = true
				break
			}
		}

		if !isIn {
			return ErrStringValueOutOfList, nil
		}
	}

	return nil, nil
}
