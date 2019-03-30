package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

func authInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(&config.Auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshaling auth config: %v", err)
		return
	}
	w.Write(bytes)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/signin/register", config.Auth.Okta.BaseURL)
	http.Redirect(w, r, url, 301)
}

func verifyToken(tokenStr string) (*jwtverifier.Jwt, error) {

	toValidate := map[string]string{}
	toValidate["aud"] = "api://default"
	toValidate["cid"] = config.Auth.Okta.ClientID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           config.Auth.Okta.BaseURL + "/oauth2/default",
		ClaimsToValidate: toValidate,
	}

	verifier := jwtVerifierSetup.New()
	verifier.SetLeeway(60)

	token, err := verifier.VerifyAccessToken(tokenStr)
	return token, err
}

func getAuth(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, bearerStr) {
		c, err := r.Cookie("auth")
		if err != nil {
			return "", err
		}
		return c.Value, nil
	}
	tokenStr := h[len(bearerStr):]
	return tokenStr, nil
}
