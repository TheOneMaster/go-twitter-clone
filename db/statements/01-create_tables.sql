CREATE TABLE Users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    displayName TEXT,
    profilePhoto TEXT,
    bannerPhoto TEXT,
    creationTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    password TEXT NOT NULL
);

CREATE TABLE Messages (
    id INTEGER PRIMARY KEY,
    parentID INTEGER,
    messageText TEXT,
    postTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    author INTEGER,
    FOREIGN KEY (author)
        REFERENCES Users (id)
    FOREIGN KEY (parentID)
        REFERENCES Messages(id)
);
