CREATE TABLE users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    moniker VARCHAR(255) NOT NULL,
    type INT DEFAULT 0 NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    about TEXT DEFAULT "Writing awesome contents at Reblog" NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
