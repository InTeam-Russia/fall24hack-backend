CREATE TYPE user_role AS ENUM ('ADMIN', 'USER');

CREATE TYPE poll_type AS ENUM ('FREE', 'RADIO');

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    role user_role NOT NULL,
    password_hash BYTEA NOT NULL,
    password_salt BYTEA NOT NULL,
    tg_link TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS polls (
    id BIGSERIAL PRIMARY KEY,
    text VARCHAR(300) NOT NULL,
    type poll_type NOT NULL,
    author_id BIGINT NOT NULL,
    cluster INT NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS user_answers (
    user_id BIGINT NOT NULL,
    poll_id BIGINT NOT NULL,
    text VARCHAR(500) NOT NULL,
    PRIMARY KEY (user_id, poll_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (poll_id) REFERENCES polls (id)
);

CREATE TABLE IF NOT EXISTS radio_answers (
    answer_id BIGINT NOT NULL,
    poll_id BIGINT NOT NULL,
    text VARCHAR(500) NOT NULL,
    PRIMARY KEY (answer_id, poll_id),
    FOREIGN KEY (poll_id) REFERENCES polls (id)
);
