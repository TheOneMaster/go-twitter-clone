CREATE TABLE Follows (
    userID INTEGER,
    followID INTEGER,
    followTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userID, followID)
    FOREIGN KEY (userID)
        REFERENCES Users(id)
    FOREIGN KEY (followID)
        REFERENCES Users(id)
);

CREATE TABLE Likes (
    messageID INTEGER,
    personID INTEGER,
    likeTime DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (messageID, personID)
    FOREIGN KEY (messageID)
        REFERENCES Messages(id)
    FOREIGN KEY (personID)
        REFERENCES Users(id)
);
