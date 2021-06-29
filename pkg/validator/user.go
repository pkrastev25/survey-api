package validator

import (
	"errors"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func RulesFirstName(firstName *string, additionalRules ...validation.Rule) *validation.FieldRules {
	rules := []validation.Rule{validation.Length(2, 20)}
	if len(additionalRules) > 0 {
		rules = append(rules, additionalRules...)
	}

	return validation.Field(firstName, rules...)
}

func RulesUserName(userName *string, additionalRules ...validation.Rule) *validation.FieldRules {
	rules := []validation.Rule{is.Alphanumeric, validation.Length(3, 20)}
	if len(additionalRules) > 0 {
		rules = append(rules, additionalRules...)
	}

	return validation.Field(userName, rules...)
}

func RulesAvatarUrl(avatarUrl *string, additionalRules ...validation.Rule) *validation.FieldRules {
	rules := []validation.Rule{is.URL}
	if len(additionalRules) > 0 {
		rules = append(rules, additionalRules...)
	}

	return validation.Field(avatarUrl, rules...)
}

func RulesPassword(password *string, additionalRules ...validation.Rule) *validation.FieldRules {
	rules := []validation.Rule{validation.Length(7, 30), validation.By(validatePassword)}
	if len(additionalRules) > 0 {
		rules = append(rules, additionalRules...)
	}

	return validation.Field(password, rules...)
}

func validatePassword(value interface{}) error {
	password, ok := value.(string)
	if !ok {
		return errors.New("")
	}

	var (
		hasUpperCase = false
		hasLowerCase = false
		hasSpecial   = false
		hasNumber    = false
	)

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
			continue
		}

		if unicode.IsUpper(char) {
			hasUpperCase = true
			continue
		}

		if unicode.IsNumber(char) {
			hasNumber = true
			continue
		}

		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			continue
		}

		if hasLowerCase && hasUpperCase && hasNumber && hasSpecial {
			return nil
		}
	}

	if !hasLowerCase {
		return errors.New("")
	}

	if !hasUpperCase {
		return errors.New("")
	}

	if !hasNumber {
		return errors.New("")
	}

	if !hasSpecial {
		return errors.New("")
	}

	return nil
}
