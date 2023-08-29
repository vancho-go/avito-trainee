CREATE TABLE IF NOT EXISTS userstosegments(
    id SERIAL PRIMARY KEY,
    user_id VARCHAR REFERENCES users(user_id),
    segment_name VARCHAR REFERENCES segments(name),
    joined TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted TIMESTAMP,
    UNIQUE (user_id, segment_name)
);