USE testdb;

CREATE TABLE users (
    Id int NOT NULL AUTO_INCREMENT,
    Username varchar(255) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    Session_Id varchar(255) NOT NULL DEFAULT 0,
    PRIMARY KEY(Id)
);

CREATE TABLE characters (
    Users_Id int NOT NULL AUTO_INCREMENT,
    Position_X int NOT NULL,
    Position_Y int NOT NULL,
    Picture_Path varchar(255) NOT NULL DEFAULT '',
    FOREIGN KEY (Users_Id) REFERENCES users(ID),
    PRIMARY KEY (Users_Id)
);
