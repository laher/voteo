package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
