
CREATE TABLE users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    moniker VARCHAR(255) NOT NULL,
    type INT DEFAULT 0 NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    about TEXT DEFAULT "Writing awesome contents at Reblog" NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE collaborator_tokens
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(225) NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE UNIQUE INDEX collaborator_tokens_email_uindex ON collaborator_tokens (email);
CREATE UNIQUE INDEX collaborator_tokens_token_uindex ON collaborator_tokens (token);

CREATE TABLE posts
(
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL,
    content TEXT NOT NULL,
    status INTEGER DEFAULT 0 NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    user_id INTEGER NOT NULL
);

CREATE UNIQUE INDEX posts_slug_uindex ON posts (slug);
