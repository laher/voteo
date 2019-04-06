package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
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
	loadConfig()
	cfg := &aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("us-east-1"),
	}
	h := newHandler(cfg)
	http.HandleFunc("/", h.templateHandler)
	http.HandleFunc("/videos", h.videosHandler)
	http.HandleFunc("/vote", h.voteHandler)
	http.HandleFunc("/yt/data", h.ytMetadataProxy)
	http.HandleFunc("/auth/settings", authInfoHandler)
	http.HandleFunc("/register", registrationHandler)
	log.Println("Listening...")
	if config.SSL {
		log.Fatal(http.Serve(autocert.NewListener(config.Address), nil))
	} else {
		log.Fatal(http.ListenAndServe(config.Address, nil))
	}
}

type video struct {
	ID    string `json:"id" dynamodbav:"id"`
	Title string `json:"title" dynamodbav:"title"`
	Votes int    `json:"votes"`
}

type vote struct {
	VideoID    string `json:"videoId" dynamodbav:"videoId"`
	PersonID   string `json:"personId,omitempty" dynamodbav:"personId"`
	PersonHash string `json:"personHash"`
	Up         bool   `json:"up" dynamodbav:"up"`
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
func newHandler(awsConfig *aws.Config) *handler {
	return &handler{db: newDynamodb(awsConfig)}
}

type handler struct {
	db *ddb
}

func (h *handler) voteHandler(w http.ResponseWriter, r *http.Request) {

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
			votes, err := h.db.getVotes()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Couldnt decode: %v", err)
				return
			}
			for _, v := range votes {
				if v.VideoID == myVote.VideoID && v.PersonID == myVote.PersonID {
					found = true
				} else {
					newVotes = append(newVotes, v)
				}
				h.db.putVote(v)
			}
			if found {
				log.Printf("Updated: %+v", myVote)
			}
		case http.MethodPost:
			// replace with received blob
			myVote := &vote{}
			d := json.NewDecoder(r.Body)
			if err := d.Decode(myVote); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("Couldnt decode: %v", err)
				return
			}
			if err := h.db.putVote(myVote); err != nil {
				respond(w, http.StatusInternalServerError)
				log.Printf("Could not put vote: %s", err)
				return
			}
		}
	}
	votes, err := h.db.getVotes()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	for _, vote := range votes {
		if personID == "unknown" || vote.PersonID != personID {
			vote.PersonID = ""
		}
	}
	tmpl, err := h.getTemplates(personID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	bs := bytes.NewBufferString("")
	videos, err := h.db.getVideos()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	doTemplate(bs, tmpl, "items.tpl", videos, votes, personID)
	itemsHTML := bs.String()
	j := struct {
		Votes     []*vote
		ItemsHTML string
	}{
		Votes:     votes,
		ItemsHTML: itemsHTML,
	}
	b, _ := json.Marshal(j)
	w.Write(b)
	log.Printf("Returned: %+v", j)
}

const bearerStr = "Bearer "

func (h *handler) videosHandler(w http.ResponseWriter, r *http.Request) {
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
			for _, v := range myVideos {
				err := h.db.put(v)
				if err != nil {
					log.Printf("Couldnt decode: %v", err)
					respond(w, http.StatusInternalServerError)
					return
				}
			}
			log.Printf("Updated: %+v", myVideos)
		}
	}
	videos, err := h.db.getVideos()
	if err != nil {
		// could not connect
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	b, _ := json.Marshal(videos)
	w.Write(b)
	log.Printf("Returned: %+v", videos)
}

func (h *handler) getTemplates(personID string) (*template.Template, error) {
	dir := "templates"
	paths := []string{
		filepath.Join(dir, "index.tpl"),
		filepath.Join(dir, "items.tpl"),
	}
	tmpl, err := template.New("index.tpl").Funcs(template.FuncMap{
		"rand": rand.Float64,
		"countVotes": func(id string) int {
			count := 0
			votes, err := h.db.getVotes()
			if err != nil {
				return 0
			}
			for _, v := range votes {
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
			votes, err := h.db.getVotes()
			if err != nil {
				return false
			}
			for _, v := range votes {
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
			votes, err := h.db.getVotes()
			if err != nil {
				return false
			}
			for _, v := range votes {
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
			votes, err := h.db.getVotes()
			if err != nil {
				return false
			}
			for _, v := range votes {
				if v.VideoID == id && v.PersonID == personID && !v.Up {
					return true
				}
			}
			return false
		},
	}).ParseFiles(paths...)
	if err != nil {
		log.Printf("Error loading templates: %s", err)
		return nil, err
	}
	return tmpl, nil
}

func (h *handler) templateHandler(w http.ResponseWriter, r *http.Request) {
	personID, err := parseAuth(r)
	if err != nil {
		// not logged in
	}

	tmpl, err := h.getTemplates(personID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
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
	videos, err := h.db.getVideos()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	votes, err := h.db.getVotes()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	doTemplate(w, tmpl, name, videos, votes, personID)
}

func doTemplate(w io.Writer, tmpl *template.Template, name string, videos []*video, votes []*vote, personID string) {
	sortByVotes(videos, votes)
	err := tmpl.Lookup(name).Execute(w, struct {
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

func (h *handler) ytMetadataProxy(w http.ResponseWriter, r *http.Request) {
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
