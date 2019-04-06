package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
)

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
}

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
