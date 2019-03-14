package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"golang.org/x/crypto/acme/autocert"
)

var (
	// some videos for the first run (when there's no db yet)
	initVideos = []*video{
		{ID: "JntjzuI5rGM", Title: "Dave Grohl Tells ...", Votes: 0},
		{ID: "X7hFERntlog", Title: "Fearless Organization", Votes: 0},
		{ID: "d_HHnEROy_w", Title: "Stop managing, start ...", Votes: -1},
		{ID: "BCkCvay4-DQ", Title: "Push Kick", Votes: 1},
	}
	config conf
)

func main() {
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/", fs)
	http.HandleFunc("/videos", videosHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/yt/data", ytDataHandler)
	http.HandleFunc("/auth/callback/login", loginCallbackHandler)
	http.HandleFunc("/auth/callback/logout", logoutCallbackHandler)
	http.HandleFunc("/auth/settings", authInfoHandler)
	loadConfig()
	loadVideos()
	loadVotes()
	log.Println("Listening...")
	if config.SSL {
		log.Fatal(http.Serve(autocert.NewListener(config.Address), nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Address, nil))
	}
}

type video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Votes int    `json:"votes"`
}

type vote struct {
	VideoID  string `json:"videoId"`
	PersonID string `json:"personId"`
	Up       bool   `json:"up"`
}

type conf struct {
	SSL     bool     `json:"ssl"`
	Address string   `json:"address"`
	Auth    authConf `json:"auth"`
}

type authConf struct {
	Type    string   `json:"type"`
	Okta    oktaConf `json:"okta"`
	Env     string   `json:"env"`
	Address string   `json:"address"`
}

type oktaConf struct {
	BaseURL     string                   `json:"baseUrl"`
	ClientID    string                   `json:"clientId"`
	RedirectURI string                   `json:"redirectUri"`
	AuthParams  map[string]interface{}   `json:"authParams"`
	idps        []map[string]interface{} `json:"idps"`
}

func loadConfig() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("couldnt read config: %v", err)
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Couldnt decode config: %v", err)
	}
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

func voteHandler(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, bearerStr) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No auth header: [%s]", h)
		return
	}
	tokenStr := h[len(bearerStr):]
	tok, err := verifyToken(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Verification failure: %v", err)
		return
	}
	claims := tok.Claims
	log.Printf("claims: %v", claims)

	switch r.Method {
	case http.MethodGet:
		// just get all
	case http.MethodDelete:
		// delete a vote
		// replace with received blob
		myVote := &vote{}
		d := json.NewDecoder(r.Body)
		err = d.Decode(myVote)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldnt decode: %v", err)
			return
		}
		//		myVote.PersonID = claims[""]
		found := false
		newVotes := []*vote{}
		for _, v := range getVotes() {
			if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
				found = true
			} else {
				newVotes = append(newVotes, v)
			}
		}
		if found {
			writeVotes(newVotes)
			log.Printf("Updated: %+v", myVote)
		}
	case http.MethodPost:
		// replace with received blob
		myVote := &vote{}
		d := json.NewDecoder(r.Body)
		err = d.Decode(myVote)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldnt decode: %v", err)
			return
		}
		//		myVote.PersonID = claims[""]
		found := false
		votes := getVotes()
		for _, v := range votes {
			if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
				v.Up = myVote.Up
				found = true
			}
		}
		if !found {
			votes = append(votes, myVote)
		}
		writeVotes(votes)
		log.Printf("Updated: %+v", myVote)
	}
	b, _ := json.Marshal(votes)
	w.Write(b)
	log.Printf("Returned: %+v", votes)
}

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

func videosHandler(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, bearerStr) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("No auth header: [%s]", h)
		return
	}
	tokenStr := h[len(bearerStr):]
	tok, err := verifyToken(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Verification failure: %v", err)
		return
	}
	claims := tok.Claims
	log.Printf("claims: %s", claims)
	switch r.Method {
	case http.MethodGet:
		// just get all
	case http.MethodPost:
		// apply a vote
	case http.MethodPut:
		// replace with received blob
		myVideos := []*video{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&myVideos)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldnt decode: %v", err)
			return
		}
		writeVideos(myVideos)
		log.Printf("Updated: %+v", myVideos)
	}
	videos := getVideos()
	b, _ := json.Marshal(videos)
	w.Write(b)
	log.Printf("Returned: %+v", videos)
}

func ytDataHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// no video id
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("No video ID")
		return
	}
	url := "https://www.youtube.com/oembed?url=http%3A//youtube.com/watch%3Fv%3D" + id
	resp, err := http.Get(url)
	if err != nil {
		// could not connect
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	l := resp.Header.Get("Content-Length")
	if l != "" {
		w.Header().Set("Content-Length", l)
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Printf("Retrieved video metadata for: %+v", id)
	} else {
		log.Printf("Error retrieving video metadata for %+v: %v", id, resp.Status)
	}
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("login redirect ... %v", r.URL.Query())
}

func logoutCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("logout redirect ... %v", r.URL.Query())
}
