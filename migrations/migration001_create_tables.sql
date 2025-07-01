CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4(),
    username VARCHAR(64) UNIQUE NOT NULL,
    nickname VARCHAR(64),
    email VARCHAR(128) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY(id)
);

CREATE TABLE conversations (
    id UUID DEFAULT uuid_generate_v4(),
    name VARCHAR(64) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT conversations_pk PRIMARY KEY(id)
);

CREATE TABLE in_conversation (
    user_id UUID,
    conversation_id UUID,
    CONSTRAINT in_conversation_pk PRIMARY KEY(user_id, conversation_id),
    CONSTRAINT in_conversation_fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT in_conversation_fk_conversation FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE TABLE messages (
    id UUID DEFAULT uuid_generate_v4(),
    send_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    body VARCHAR(256) NOT NULL,
    sender_id UUID NOT NULL,
    conversation_id UUID NOT NULL,
    CONSTRAINT messages_pk PRIMARY KEY(id),
    CONSTRAINT messages_fk_sender FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT messages_fk_conversation FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);

CREATE TABLE is_friend (
    user1_id UUID,
    user2_id UUID,
    CONSTRAINT is_friend_pk PRIMARY KEY(user1_id, user2_id),
    CONSTRAINT is_friend_fk_user1 FOREIGN KEY(user1_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT is_friend_fk_user2 FOREIGN KEY(user2_id) REFERENCES users(id) ON DELETE CASCADE
);
