package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

const bearerStr = "Bearer "

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

func parseAuth(r *http.Request) (string, error) {
	tokenStr, err := getAuth(r)
	if err != nil {
		return "", err
	}
	tok, err := verifyToken(tokenStr)
	if err != nil {
		return "", err
	}
	claims := tok.Claims
	log.Printf("claims: %v", claims)
	personIDI, ok := claims["sub"]
	if !ok {
		return "", errors.New("claims 'sub' field does not exist")
	}
	personID, ok := personIDI.(string)
	if !ok {
		return "", errors.New("invalid claims 'sub' field")
	}
	return personID, nil
}
