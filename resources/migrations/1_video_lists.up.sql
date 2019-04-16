set timezone='UTC';

CREATE TABLE IF NOT EXISTS video_lists (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    creator VARCHAR(255) NOT NULL,
    inserted_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX video_lists_creator ON video_lists (creator);
CREATE INDEX video_lists_inserted_at ON video_lists (inserted_at);


CREATE type VIDEO_SOURCE AS ENUM('youtube');

CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL,
    source VIDEO_SOURCE NOT NULL,
    video_list_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    creator VARCHAR(255) NOT NULL,
    inserted_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX videos_creator ON videos (creator);
CREATE INDEX videos_video_list_id ON videos (video_list_id);
CREATE INDEX videos_inserted_at ON videos (inserted_at);


CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    video_id INT NOT NULL,
    video_list_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    creator VARCHAR(255) NOT NULL,
    up bool NOT NULL,
    inserted_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP with time zone default CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX votes_creator ON votes (creator);
CREATE INDEX votes_video_list_id ON votes (video_list_id);
CREATE INDEX votes_video_id ON votes (video_id);
CREATE INDEX votes_inserted_at ON votes (inserted_at);

CREATE UNIQUE INDEX votes_video ON votes (video_id, creator);
