package main

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	row := db.db.QueryRowx(`SELECT * FROM users WHERE id = $1 ORDER BY rank`, id)
	f := &user{}
	if err := row.StructScan(f); err != nil {
		return nil, err
	}
	return f, nil
}

//func (db *pdb) createUser(u *user) error {}

func (db *pdb) getVideoList(id int) (*videoList, error) {
	row := db.db.QueryRowx(`SELECT * FROM videoLists WHERE id = $1 ORDER BY rank`, id)
	f := &videoList{}
	if err := row.StructScan(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (db *pdb) getVideoLists(userID string) ([]*videoList, error) {
	rows, err := db.db.Queryx(`SELECT * FROM video_lists WHERE creator_id = $1 ORDER by inserted_at`, userID)
	if err != nil {
		return nil, err
	}

	videoLists := []*videoList{}
	for rows.Next() {
		f := &videoList{}
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
	sqlUpdate := `INSERT INTO videos (video_list_id, video, creator, inserted_at)
		  VALUES ($1, $2, $3, NOW())`
	res, err := db.db.Exec(sqlUpdate,
		v.VideoListID, v.ID, v.CreatorID)
	if err != nil {
		return errors.Wrap(err, "Failed to insert video")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Failed to get rows affected")
	}
	if rowsAffected != 1 {
		err = sql.ErrNoRows
		return err
	}
	return nil
}

func (db *pdb) putVote(v *vote) error {
	sqlUpdate := `INSERT INTO videos (video_list_id, video, creator_id, up, inserted_at)
		  VALUES ($1, $2, $3, NOW()) RETURNING id`
	res, err := db.db.Exec(sqlUpdate,
		v.VideoListID, v.VideoID, v.Up, v.PersonID)
	if err != nil {
		return errors.Wrap(err, "Failed to insert video")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Failed to get rows affected")
	}
	v.ID = int(id)
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
	sqlUpdate := `INSERT INTO video_lists (creator_id, title, inserted_at)
		  VALUES ($1,$2,NOW()) RETURNING id`
	res, err := db.db.Exec(sqlUpdate,
		vl.CreatorID, vl.Title)
	if err != nil {
		return errors.Wrap(err, "Failed to insert video_list")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Failed to get rows affected")
	}
	vl.ID = int(id)
	return nil
}
