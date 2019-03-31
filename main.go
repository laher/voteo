package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

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
	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", fs)
	http.HandleFunc("/", templateHandler)
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

	personID, err := parseAuth(r)
	switch r.Method {
	case http.MethodGet:
		// no auth needed
		if err != nil {
			personID = "unknown"
		}
	default:
		if err != nil {
			respond(w, http.StatusUnauthorized)
			log.Printf("auth failure: %v", err)
			return
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
		_, err := parseAuth(r)
		if err != nil {
			log.Printf("Auth failure: %v", err)
			respond(w, http.StatusUnauthorized)
			return
		}
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

func templateHandler(w http.ResponseWriter, r *http.Request) {
	personID, err := parseAuth(r)
	if err != nil {
		// not logged in
	}
	dir := "templates"
	paths := []string{
		filepath.Join(dir, "index.tpl"),
		filepath.Join(dir, "items.tpl"),
	}
	tmpl, err := template.New("index.tpl").Funcs(template.FuncMap{
		"countVotes": func(id string) int {
			count := 0
			for _, v := range db.getVotes() {
				if v.VideoID == id {
					if v.Up {
						count++
					} else {
						count--
					}
				}
			}
			return count
		},
		"haveIVoted": func(id string) bool {
			if personID == "" {
				return false
			}
			for _, v := range db.getVotes() {
				if v.VideoID == id && v.PersonID == personID {
					return true
				}
			}
			return false
		},
		"haveIUpvoted": func(id string) bool {
			if personID == "" {
				return false
			}
			for _, v := range db.getVotes() {
				if v.VideoID == id && v.PersonID == personID && v.Up {
					return true
				}
			}
			return false
		},
		"haveIDownvoted": func(id string) bool {
			if personID == "" {
				return false
			}
			for _, v := range db.getVotes() {
				if v.VideoID == id && v.PersonID == personID && !v.Up {
					return true
				}
			}
			return false
		},
	}).ParseFiles(paths...)
	if err != nil {
		log.Fatalf("loading templates: %s", err)
	}
	name := ""
	parts := strings.Split(r.URL.Path, "/")
	part := parts[len(parts)-1]
	switch part {
	case "items":
		name = "items.tpl"
	case "":
		// ok
		name = "index.tpl"
	default:
		respond(w, http.StatusNotFound)
		return
	}
	videos := db.getVideos()
	votes := db.getVotes()
	sortByVotes(videos, votes)
	err = tmpl.Lookup(name).Execute(w, struct {
		PersonID string
		Items    []*video
	}{
		PersonID: personID,
		Items:    videos,
	})
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}

func sortByVotes(videos []*video, votes []*vote) {
	// reset votes ...
	for _, video := range videos {
		video.Votes = 0
	}
	for _, vote := range votes {
		for _, video := range videos {
			if video.ID == vote.VideoID {
				if vote.Up {
					video.Votes++
				} else {

					video.Votes--
				}
			}
		}
	}
	sort.Sort(sort.Reverse(videosByVotes(videos)))
}

type videosByVotes []*video

func (a videosByVotes) Len() int           { return len(a) }
func (a videosByVotes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a videosByVotes) Less(i, j int) bool { return a[i].Votes < a[j].Votes }

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
