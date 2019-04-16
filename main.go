package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

var (
	// some videos for the first run (when there's no db yet)
	/*
		initVideos = []*video{
			{ID: "JntjzuI5rGM", Title: "Dave Grohl Tells ...", Votes: 0},
			{ID: "X7hFERntlog", Title: "Fearless Organization", Votes: 0},
			{ID: "d_HHnEROy_w", Title: "Stop managing, start ...", Votes: -1},
			{ID: "BCkCvay4-DQ", Title: "Push Kick", Votes: 1},
		}
	*/
	config conf
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", fs)
	loadConfig()
	dsn := fmt.Sprintf(config.PostgresURL, config.PostgresUser, config.PostgresPassword)
	db, err := newDB(dsn)
	if err != nil {
		logrus.Fatalf("Could not connect to DB, (%s), user: %s, error: %v", config.PostgresURL, config.PostgresUser, err)
	}
	if err = runMigrationsSource(db.db.DB); err != nil {
		logrus.WithError(err).Fatalf("Could not run migrations (%s), user: %s,", config.PostgresURL, config.PostgresUser)
	}
	logrus.Info("Migrations complete")
	h := newHandler(db)
	http.HandleFunc("/", h.templateHandler)
	http.HandleFunc("/videos", h.videosHandler)
	http.HandleFunc("/vote", h.voteHandler)
	http.HandleFunc("/videoLists", h.videoListsHandler)
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
	SSL              bool     `json:"ssl"`
	Address          string   `json:"address"`
	Auth             authConf `json:"auth"`
	PostgresURL      string   `json:"postgresURL"`
	PostgresUser     string   `json:"postgresUser"`
	PostgresPassword string   `json:"postgresPassword"`
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

func doTemplate(w io.Writer, tmpl *template.Template, name string, myVideoList *videoList, myVideoLists []*videoList, personID string) error {
	//sortByVotes(videos, votes)
	err := tmpl.Lookup(name).Execute(w, struct {
		PersonID   string
		VideoList  *videoList
		VideoLists []*videoList
	}{
		PersonID:   personID,
		VideoList:  myVideoList,
		VideoLists: myVideoLists,
	})
	return err
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(`{ "error": "` + http.StatusText(statusCode) + `" }`)
	w.Write(b)
}
