package handlers

import (
	"net/http"
	"strings"
	"Gopchat/tokens"
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

func authenticate(r *http.Request)(int, bool){
	token := getCredentials(r)

	if token ==""{
		return 0, false
	}

	userId, ok:= tokens.GetUser(token)

	return userId, ok
}

