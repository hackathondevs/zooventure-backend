CREATE TABLE Users (
    ID BIGINT UNSIGNED AUTO_INCREMENT,

    Email VARCHAR(255) NOT NULL UNIQUE,
    Password CHAR(60) NOT NULL,
    Name VARCHAR(255) NOT NULL UNIQUE,
    ProfilePicture TEXT NOT NULL DEFAULT '',
    Active BOOLEAN NOT NULL DEFAULT FALSE,
    Admin BOOLEAN NOT NULL DEFAULT FALSE,
    Premium BOOLEAN NOT NULL DEFAULT FALSE,
    
    CreatedAt DATETIME NOT NULL DEFAULT NOW(),
    UpdatedAt DATETIME NOT NULL ON UPDATE NOW(),
    
    PRIMARY KEY (ID)
) ENGINE = INNODB DEFAULT CHARSET = UTF8;

INSERT INTO Users (
    Email, 
    Password, 
    Name,
    Active,
    Admin,
    Premium
)
VALUE (
    "otter.whopper@gmail.com", 
    "12345678", 
    "Otter Whopper",
    TRUE,
    TRUE,
    TRUE
);