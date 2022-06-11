package hw09structvalidator

import (
	"errors"
	"strconv"
	"strings"
)

var ErrIntValueTooSmall = errors.New("value too small")

var ErrIntValueTooLarge = errors.New("value too large")

var ErrIntValueOutOfList = errors.New("value out of list")

func validateIntRules(field string, value int, rules validationRules) (ValidationErrors, error) {
	errors := make(ValidationErrors, 0)
	for _, rule := range rules {
		validateErr, err := validateInt(value, rule)
		if err != nil {
			return nil, err
		}

		if validateErr != nil {
			errors = append(errors, ValidationError{Field: field, Err: validateErr})
		}
	}

	return errors, nil
}

func validateInt(value int, rule validationRule) (error, error) {
	switch rule.Name {
	case "min":
		min, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}

		if value < min {
			return ErrIntValueTooSmall, nil
		}
	case "max":
		max, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}

		if value > max {
			return ErrIntValueTooLarge, nil
		}
	case "in":
		values := strings.Split(rule.Value, ",")
		isIn := false
		for _, v := range values {
			intValue, err := strconv.Atoi(v)
			if err == nil && value == intValue {
				isIn = true
				break
			}
		}

		if !isIn {
			return ErrIntValueOutOfList, nil
		}
	}

	return nil, nil
}
