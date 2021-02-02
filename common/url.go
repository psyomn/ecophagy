package common

import (
	"errors"
	"strings"
)

func PartsOfURLSafe(url string) ([]string, error) {
	if strings.Contains(url, `..`) {
		return nil, errors.New("unsafe url: detected .. attempt")
	}

	parts := strings.Split(url, "/")
	ret := []string{}

	for index := range parts {
		if parts[index] == "" {
			continue
		}

		ret = append(ret, parts[index])
	}

	return ret, nil
}
