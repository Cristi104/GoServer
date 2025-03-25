CREATE TABLE users (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) UNIQUE,
    email VARCHAR(128) UNIQUE,
    password VARCHAR(256)
);

CREATE TABLE conversations (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    create_date DATETIME NOT NULL
);

CREATE TABLE in_conversation (
    user_id INTEGER,
    conversation_id INTEGER,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE TABLE messages (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    send_date DATETIME,
    body VARCHAR(256),
    sender_id INTEGER,
    conversation_id INTEGER,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);
