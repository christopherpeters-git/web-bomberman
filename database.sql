USE testdb;

CREATE TABLE users (
    Id int NOT NULL AUTO_INCREMENT,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    Session_Id varchar(255) NOT NULL DEFAULT 0,
    Picture_Path varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY(Id)
);
