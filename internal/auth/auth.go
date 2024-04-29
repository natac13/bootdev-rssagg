package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("missing api key")
	}

	vals := strings.Split(apiKey, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid api key")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid api key")
	}

	return vals[1], nil
}
