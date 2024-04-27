CREATE TABLE Merchants (
    ID BIGINT UNSIGNED AUTO_INCREMENT,
    
    Name VARCHAR(255) NOT NULL,
    Code VARCHAR(255) NOT NULL UNIQUE,

    CreatedAt DATETIME NOT NULL DEFAULT NOW(),
    UpdatedAt DATETIME NOT NULL ON UPDATE NOW() DEFAULT NOW(),
    
    PRIMARY KEY (ID)
) ENGINE = INNODB DEFAULT CHARSET = UTF8;

INSERT INTO Merchants (
    Name,
    Code
) VALUES (
    "Stand Teh Poci",
    "KBS-001"
),
(
    "Stand Kentang Goreng",
    "KBS-002"
);