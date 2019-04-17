package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type pdb struct {
	db *sqlx.DB
}

func newDB(url string) (*pdb, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open database")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "could not ping database")
	}
	db.SetMaxOpenConns(30)
	db.SetConnMaxLifetime(3 * time.Minute)
	return &pdb{db: db}, nil
}

//func (db *pdb) createTables() error {}
//func (db *pdb) dropTables() error   {}

func (db *pdb) getUser(id string) (*user, error) {
	row := db.db.QueryRowx(`SELECT * FROM users WHERE id = $1`, id)
	f := &user{}
	if err := row.StructScan(f); err != nil {
		return nil, err
	}
	return f, nil
}

//func (db *pdb) createUser(u *user) error {}

func (db *pdb) getVideoList(id int) (*videoList, error) {
	row := db.db.QueryRowx(`SELECT * FROM video_lists WHERE id = $1`, id)
	f := &videoList{}
	if err := row.StructScan(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (db *pdb) getVideoLists(userID string) ([]*videoList, error) {
	rows, err := db.db.Queryx(`SELECT * FROM video_lists WHERE creator = $1 ORDER by inserted_at`, userID)
	if err != nil {
		return nil, err
	}

	videoLists := []*videoList{}
	for rows.Next() {
		f := &videoList{}
		f.Videos = []*video{}
		f.Votes = []*vote{}
		if err := rows.StructScan(f); err != nil {
			return nil, err
		}

		videoLists = append(videoLists, f)
	}
	return videoLists, nil
}

func (db *pdb) getVideosForList(videoListID int) ([]*video, error) {
	rows, err := db.db.Queryx(`SELECT * FROM videos WHERE video_list_id = $1 ORDER by inserted_at`, videoListID)
	if err != nil {
		return nil, err
	}
	videos := []*video{}
	for rows.Next() {
		f := &video{}
		if err := rows.StructScan(f); err != nil {
			return nil, err
		}

		videos = append(videos, f)
	}
	return videos, nil
}

func (db *pdb) addVideoToList(v *video) error {
	sqlUpdate := `INSERT INTO videos (title, video_list_id, source_id, source, creator, inserted_at, updated_at)
		  VALUES ($1, $2, $3, 'youtube', $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
	err := db.db.QueryRowx(sqlUpdate,
		v.Title, v.VideoListID, v.SourceID, v.CreatorID).Scan(&v.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to insert video")
	}
	return nil
}

func (db *pdb) updateVote(v *vote) error {
	sqlUpdate := `UPDATE votes SET up=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`
	res, err := db.db.Exec(sqlUpdate,
		v.VideoID, v.Up)
	if err != nil {
		return errors.Wrap(err, "Failed to update vote")
	}
	if rows, err := res.RowsAffected(); true {
		if err != nil {
			return err
		}
		if rows != 1 {
			return errors.New("Failed to update vote")
		}
	}
	return nil
}

func (db *pdb) deleteVote(v *vote) error {
	sqlUpdate := `DELETE FROM votes WHERE id=$1`
	res, err := db.db.Exec(sqlUpdate,
		v.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to delete vote")
	}
	if rows, err := res.RowsAffected(); true {
		if err != nil {
			return err
		}
		if rows != 1 {
			return errors.New(fmt.Sprintf("Failed to delete vote (%d rows affected)", rows))
		}
	}
	return nil
}

func (db *pdb) createVote(v *vote) error {
	sqlUpdate := `INSERT INTO votes (video_list_id, video_id, up, creator, inserted_at, updated_at)
		  VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
	err := db.db.QueryRowx(sqlUpdate,
		v.VideoListID, v.VideoID, v.Up, v.PersonID).Scan(&v.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to insert vote")
	}
	return nil
}

func (db *pdb) getVotes(videoListID int) ([]*vote, error) {
	rows, err := db.db.Queryx(`SELECT * FROM votes WHERE video_list_id = $1 ORDER by inserted_at`, videoListID)
	if err != nil {
		return nil, err
	}
	v := []*vote{}
	for rows.Next() {
		f := &vote{}
		if err := rows.StructScan(f); err != nil {
			return nil, err
		}

		v = append(v, f)
	}
	return v, nil
}

func (db *pdb) createVideoList(vl *videoList) error {
	sqlUpdate := `INSERT INTO video_lists (creator, title, inserted_at, updated_at)
		  VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
	err := db.db.QueryRowx(sqlUpdate,
		vl.CreatorID, vl.Title).Scan(&vl.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to insert video_list")
	}
	return nil
}
