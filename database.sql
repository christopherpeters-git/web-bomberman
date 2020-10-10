USE testdb;

CREATE TABLE users (
    Id int NOT NULL AUTO_INCREMENT,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    Session_Id varchar(255) UNIQUE,
    Games_won int,
    PRIMARY KEY(Id)
);
