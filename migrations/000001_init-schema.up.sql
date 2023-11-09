CREATE TABLE IF NOT EXISTS wisdoms
(
    id         UUID                        NOT NULL CHECK (id != '00000000-0000-0000-0000-000000000000'),
    value      TEXT                        NOT NULL,
    author_id  UUID                        NOT NULL CHECK (author_id != '00000000-0000-0000-0000-000000000000'),

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS authors
(
    id         UUID                        NOT NULL CHECK (id != '00000000-0000-0000-0000-000000000000'),
    name       TEXT                        NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,

    PRIMARY KEY (id)
);
