CREATE TABLE Users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    displayName TEXT,
    photo TEXT,
    creationTime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Messages (
    id INTEGER PRIMARY KEY,
    messageText TEXT,
    time DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    author INTEGER,
    FOREIGN KEY (author)
        REFERENCES Users (id)
);
