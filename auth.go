package main

import (
	"encoding/json"
	"log"
	"net/http"

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
