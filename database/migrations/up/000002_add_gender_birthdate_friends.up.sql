ALTER TABLE users
    ADD COLUMN gender     VARCHAR(10)  DEFAULT NULL,
    ADD COLUMN birth_date DATE         DEFAULT NULL;

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,

    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT no_self_friend CHECK (user_id != friend_id)
);
