package hw09structvalidator

import (
	"strings"
)

type validationRule struct {
	Name  string
	Value string
}

type validationRules []validationRule

func parseRules(stringRules string) validationRules {
	splitedRules := strings.Split(stringRules, "|")
	rules := make(validationRules, 0)

	for _, value := range splitedRules {
		rules = append(rules, parseSingleRule(value))
	}

	return rules
}

func parseSingleRule(rule string) validationRule {
	splitedRule := strings.SplitN(rule, ":", 2)
	name := splitedRule[0]

	var value string
	if len(splitedRule) > 1 {
		value = splitedRule[1]
	}

	return validationRule{Name: name, Value: value}
}
