package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	http.HandleFunc("/yt/data", ytMetadataProxy)
	http.HandleFunc("/auth/settings", authInfoHandler)
	http.HandleFunc("/register", registrationHandler)
	loadConfig()
	db.loadVideos()
	db.loadVotes()
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
	VideoID    string `json:"videoId"`
	PersonID   string `json:"personId,omitempty"`
	PersonHash string `json:"personHash"`
	Up         bool   `json:"up"`
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

func voteHandler(w http.ResponseWriter, r *http.Request) {

	var personID = "unknown"
	switch r.Method {
	case http.MethodGet:
		// no auth needed
	default:
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, bearerStr) {
			respond(w, http.StatusUnauthorized)
			log.Printf("No auth header: [%s]", h)
			return
		}
		tokenStr := h[len(bearerStr):]
		tok, err := verifyToken(tokenStr)
		if err != nil {
			respond(w, http.StatusUnauthorized)
			log.Printf("Verification failure: %v", err)
			return
		}
		claims := tok.Claims
		log.Printf("claims: %v", claims)

		personIDI, ok := claims["sub"]
		if ok {
			personID, ok = personIDI.(string)
		}
		switch r.Method {
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
			hash := sha256.New()
			hash.Write([]byte(personID))
			myVote.PersonHash = fmt.Sprintf("%x", hash.Sum(nil))
			found := false
			newVotes := []*vote{}
			for _, v := range db.getVotes() {
				if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
					found = true
				} else {
					newVotes = append(newVotes, v)
				}
			}
			if found {
				db.writeVotes(newVotes)
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
			found := false
			myVotes := db.getVotes()
			for _, v := range myVotes {
				if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
					v.Up = myVote.Up
					found = true
				}
			}
			if !found {
				myVotes = append(myVotes, myVote)
			}
			db.writeVotes(myVotes)
			log.Printf("Updated: %+v", myVote)
		}
	}
	votes := db.getVotes()
	for _, vote := range votes {
		if personID == "unknown" || vote.PersonID != personID {
			vote.PersonID = ""
		}
	}
	b, _ := json.Marshal(votes)
	w.Write(b)
	log.Printf("Returned: %+v", votes)
}

const bearerStr = "Bearer "

func videosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// no auth required ...
	default:
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, bearerStr) {
			log.Printf("No auth header: [%s]", h)
			respond(w, http.StatusUnauthorized)
			return
		}
		tokenStr := h[len(bearerStr):]
		tok, err := verifyToken(tokenStr)
		if err != nil {
			log.Printf("Verification failure: %v", err)
			respond(w, http.StatusUnauthorized)
			return
		}
		claims := tok.Claims
		log.Printf("claims: %s", claims)
		switch r.Method {
		case http.MethodPost:
			// apply a vote
		case http.MethodPut:
			// replace with received blob
			myVideos := []*video{}
			d := json.NewDecoder(r.Body)
			err := d.Decode(&myVideos)
			if err != nil {
				log.Printf("Couldnt decode: %v", err)
				respond(w, http.StatusBadRequest)
				return
			}
			db.writeVideos(myVideos)
			log.Printf("Updated: %+v", myVideos)
		}
	}
	videos := db.getVideos()
	b, _ := json.Marshal(videos)
	w.Write(b)
	log.Printf("Returned: %+v", videos)
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(`{ "error": "` + http.StatusText(statusCode) + `" }`)
	w.Write(b)
}

func ytMetadataProxy(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// no video id
		respond(w, http.StatusBadRequest)
		log.Printf("No video ID")
		return
	}
	url := "https://www.youtube.com/oembed?url=http%3A//youtube.com/watch%3Fv%3D" + id
	resp, err := http.Get(url)
	if err != nil {
		// could not connect
		respond(w, http.StatusInternalServerError)
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
