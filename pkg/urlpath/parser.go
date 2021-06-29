package urlpath

import (
	"errors"
	"regexp"
	"strings"
)

const (
	regexAny = "(.*)"
)

type UrlParser struct {
}

func NewUrlParser() UrlParser {
	return UrlParser{}
}

func (parser UrlParser) ParseParams(path string, recipe string) (map[string]string, error) {
	var keys []string
	recipeWithRegex := recipe
	params := make(map[string]string)

	for true {
		start := strings.Index(recipeWithRegex, "{")
		end := strings.Index(recipeWithRegex, "}")
		if start == -1 && end == -1 {
			break
		}

		if start == -1 || end == -1 {
			return nil, errors.New("")
		}

		key := recipeWithRegex[start+1 : end]
		keys = append(keys, key)
		recipeWithRegex = strings.ReplaceAll(recipeWithRegex, "{"+key+"}", regexAny)
	}

	if len(keys) <= 0 {
		return params, nil
	}

	values, err := parser.applyRegex(path, recipeWithRegex)
	if err != nil {
		return nil, err
	}

	if len(values) <= 0 {
		return params, nil
	}

	if len(keys) != len(values) {
		return nil, errors.New("")
	}

	for i := 0; i < len(keys); i++ {
		params[keys[i]] = values[i]
	}

	return params, nil
}

func (parser UrlParser) applyRegex(path string, recipeWithRegex string) ([]string, error) {
	regexRecipe, err := regexp.Compile(recipeWithRegex)
	if err != nil {
		return nil, err
	}

	results := regexRecipe.FindStringSubmatch(path)
	if len(results) <= 1 {
		return nil, errors.New("")
	}

	return results[1:], nil
}
