CREATE TABLE users (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64),
    email VARCHAR(128) UNIQUE,
    password VARCHAR(256)
);

CREATE TABLE messages (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    send_date DATETIME,
    body VARCHAR(256),
    sender_id INTEGER,
    receiver_id INTEGER,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE RESTRICT
);