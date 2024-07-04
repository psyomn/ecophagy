-- +goose Up

CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL UNIQUE,
        password BLOB NOT NULL,

        created_at DATETIME,
        updated_at DATETIME
);

CREATE TABLE notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        owner_id INTEGER,
        title VARCHAR NOT NULL,
        comment VARCHAR,
        contents VARCHAR,
        view_mode INTEGER NOT NULL CHECK (view_mode = 0 OR view_mode = 1),

        created_at DATETIME,
        updated_at DATETIME,

        CONSTRAINT fk_users
              FOREIGN KEY (owner_id) REFERENCES users (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE users;
DROP TABLE notes;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
