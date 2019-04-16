package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func newHandler(db *pdb) *handler {
	return &handler{db: db}
}

type handler struct {
	db *pdb
}

func (h *handler) listHandler(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) voteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		respond(w, http.StatusBadRequest)
		log.Printf("No Video List ID received: %s", err)
		return
	}
	myVideoList, err := h.db.getVideoList(id)
	if err != nil {
		// could not connect
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
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
			votes, err := h.db.getVotes(id)
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
	votes, err := h.db.getVotes(id)
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
	j := struct {
		Votes     []*vote
		VideoList *videoList
	}{
		Votes:     votes,
		VideoList: myVideoList,
	}
	b, _ := json.Marshal(j)
	w.Write(b)
}

func (h *handler) videoListsHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	switch r.Method {
	case http.MethodGet:
		// no auth required ...
		var err error
		id, err = strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			respond(w, http.StatusBadRequest)
			log.Printf("No Video List ID received: %s", err)
			return
		}
	default:
		_, err := parseAuth(r)
		if err != nil {
			log.Printf("Auth failure: %v", err)
			respond(w, http.StatusUnauthorized)
			return
		}
		switch r.Method {
		case http.MethodPut:
			myVideoList := &videoList{}
			d := json.NewDecoder(r.Body)
			err := d.Decode(&myVideoList)
			if err != nil {
				log.Printf("Couldnt decode: %v", err)
				respond(w, http.StatusBadRequest)
				return
			}
			err = h.db.createVideoList(myVideoList)
			if err != nil {
				log.Printf("Couldnt create: %v", err)
				respond(w, http.StatusInternalServerError)
				return
			}
			id = myVideoList.ID
			log.Printf("Created: %+v", myVideoList)
		}
	}
	myVideoList, err := h.db.getVideoList(id)
	if err != nil {
		// could not connect
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	b, _ := json.Marshal(myVideoList)
	w.Write(b)
	log.Printf("Returned: %+v", myVideoList)
}

func (h *handler) videosHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		respond(w, http.StatusBadRequest)
		log.Printf("No Video List ID received: %s", err)
		return
	}
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
				v.VideoListID = id
				err := h.db.addVideoToList(v)
				if err != nil {
					log.Printf("Couldnt decode: %v", err)
					respond(w, http.StatusInternalServerError)
					return
				}
			}
			log.Printf("Updated: %+v", myVideos)
		}
	}
	videos, err := h.db.getVideosForList(id)
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

func (h *handler) getTemplates(videoListID int, personID string) (*template.Template, error) {
	glob := "templates/*.tpl"
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("index.tpl").Funcs(template.FuncMap{
		"rand": rand.Float64,
		"countVotes": func(id int) int {
			count := 0
			votes, err := h.db.getVotes(videoListID)
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
		"haveIVoted": func(id int) bool {
			if personID == "" {
				return false
			}
			votes, err := h.db.getVotes(videoListID)
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
		"json": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"haveIUpvoted": func(id int) bool {
			if personID == "" {
				return false
			}
			votes, err := h.db.getVotes(videoListID)
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
		"haveIDownvoted": func(id int) bool {
			if personID == "" {
				return false
			}
			votes, err := h.db.getVotes(videoListID)
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
	idS := r.URL.Query().Get("id")
	var id int
	if idS != "" {
		id, err = strconv.Atoi(idS)
		if err != nil {
			respond(w, http.StatusBadRequest)
			log.Printf("No Video List ID received: %s", err)
			return
		}
	}
	tmpl, err := h.getTemplates(id, personID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		log.Printf("Could not fetch metadata: %s", err)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	part := parts[len(parts)-1]
	switch part {
	case "video-list":
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			respond(w, http.StatusBadRequest)
			log.Printf("No Video List ID received: %s", err)
			return
		}
		myVideoList, err := h.db.getVideoList(id)
		if err != nil {
			respond(w, http.StatusNotFound)
			log.Printf("Could not fetch metadata: %s", err)
			return
		}
		log.Printf("video list found: %d, %+v, %d", id, myVideoList, len(myVideoList.Videos))
		name := part + ".tpl"
		log.Printf("Resolved to template %s", name)
		err = tmpl.Lookup(name).Execute(w, struct {
			Page      string
			PageName  string
			PersonID  string
			VideoList *videoList
		}{
			Page:      name,
			PageName:  strings.Replace(name, ".tpl", "", -1),
			PersonID:  personID,
			VideoList: myVideoList,
		})
		if err != nil {
			respond(w, http.StatusInternalServerError)
			log.Printf("Could not fetch metadata: %s", err)
			return
		}
	case "", "index":
		myVideoLists, err := h.db.getVideoLists(personID)
		if err != nil {
			respond(w, http.StatusNotFound)
			log.Printf("Could not fetch metadata: %s", err)
			return
		}
		log.Printf("video lists found: %d", len(myVideoLists))
		name := "index.tpl"
		err = tmpl.Lookup(name).Execute(w, struct {
			Page       string
			PageName   string
			PersonID   string
			VideoLists []*videoList
		}{
			Page:       name,
			PageName:   strings.Replace(name, ".tpl", "", -1),
			PersonID:   personID,
			VideoLists: myVideoLists,
		})
		if err != nil {
			respond(w, http.StatusInternalServerError)
			log.Printf("Could not fetch metadata: %s", err)
			return
		}
	default:
		respond(w, http.StatusNotFound)
		return
	}

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
