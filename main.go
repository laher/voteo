package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	var (
		prod bool
	)
	flag.BoolVar(&prod, "prod", false, "is this prod?")
	flag.Parse()
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/", fs)
	http.HandleFunc("/videos", videosHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/yt/data", ytDataHandler)
	http.HandleFunc("/auth/callback/login", loginCallbackHandler)
	http.HandleFunc("/auth/callback/logout", logoutCallbackHandler)
	loadVideos()
	loadVotes()
	log.Println("Listening...")
	if prod {
		log.Fatal(http.Serve(autocert.NewListener("voteo.laher.net.nz"), nil))
	} else {
		log.Fatal(http.ListenAndServe("localhost:3000", nil))
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

var (
	videos     = []*video{}
	votes      = []*vote{}
	initVideos = []*video{
		{ID: "JntjzuI5rGM", Title: "Dave Grohl Tells ...", Votes: 0},
		{ID: "X7hFERntlog", Title: "Fearless Organization", Votes: 0},
		{ID: "d_HHnEROy_w", Title: "Stop managing, start ...", Votes: -1},
		{ID: "BCkCvay4-DQ", Title: "Push Kick", Votes: 1},
	}
	initVotes = []*vote{
		{VideoID: "BCkCvay4-DQ", PersonID: "am@voteo", Up: true},
	}
	lock = sync.RWMutex{}
)

func loadVideos() {
	lock.Lock()
	defer lock.Unlock()
	b, err := ioutil.ReadFile("videos.json")
	if err != nil {
		if os.IsNotExist(err) {
			videos = initVideos
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVideos := []*video{}
	err = json.Unmarshal(b, &myVideos)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	videos = myVideos
}

func loadVotes() {
	lock.Lock()
	defer lock.Unlock()
	b, err := ioutil.ReadFile("votes.json")
	if err != nil {
		if os.IsNotExist(err) {
			videos = initVideos
			return
		}
		log.Fatalf("couldnt read db: %v", err)
	}
	myVotes := []*vote{}
	err = json.Unmarshal(b, &myVotes)
	if err != nil {
		log.Fatalf("Couldnt decode db: %v", err)
	}
	votes = myVotes
}
func writeVideos() {
	lock.Lock()
	defer lock.Unlock()
	b, _ := json.MarshalIndent(videos, "", "  ")
	err := ioutil.WriteFile("videos.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
}

func writeVotes() {
	lock.Lock()
	defer lock.Unlock()
	b, _ := json.MarshalIndent(votes, "", "  ")
	err := ioutil.WriteFile("votes.json", b, 0644)
	if err != nil {
		log.Fatalf("Couldnt write db: %v", err)
	}
}

func verifyToken(tokenStr string) (*jwtverifier.Jwt, error) {

	toValidate := map[string]string{}
	toValidate["aud"] = "api://default"
	toValidate["cid"] = "0oabsbm6ga3Sy1tIf356"

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           "https://dev-343286.okta.com/oauth2/default",
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
		lock.Lock()
		found := false
		newVotes := []*vote{}
		for _, v := range votes {
			if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
				found = true
			} else {
				newVotes = append(newVotes, v)
			}
		}
		votes = newVotes
		lock.Unlock()
		if found {
			writeVotes()
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
		lock.Lock()
		found := false
		for _, v := range votes {
			if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
				v.Up = myVote.Up
				found = true
			}
		}
		if !found {
			votes = append(votes, myVote)
		}
		lock.Unlock()
		writeVotes()
		log.Printf("Updated: %+v", myVote)
	}
	lock.RLock()
	defer lock.RUnlock()
	b, _ := json.Marshal(votes)
	w.Write(b)
	log.Printf("Returned: %+v", votes)
}

const bearerStr = "Bearer "

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
		lock.Lock()
		videos = myVideos
		lock.Unlock()
		writeVideos()
		log.Printf("Updated: %+v", myVideos)
	}
	lock.RLock()
	defer lock.RUnlock()
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
