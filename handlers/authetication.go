package handlers

import (
	"net/http"
	"strings"
)

const(
	AUTH_PREFIX="Token"
	AUTH_HEADER = "Authetication"
)

func getCredentials(r *http.Request) string{
	header := r.Header.Get(AUTH_HEADER)
	token := strings.TrimPrefix(header, AUTH_PREFIX)
	return token
}

