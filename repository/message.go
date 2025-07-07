package repository

const INSERT_MESSAGE_SQL = "INSERT INTO messages(body, sender_id, conversation_id) VALUES ($1, $2, $3) RETURNING id, send_date;"
const SELECT_MESSAGES_BY_CONVERSATION_SQL = "SELECT * FROM messages WHERE conversation_id = $1;"
const UPDATE_MESSAGE_SQL = "UPDATE messages SET body = $1 WHERE id = $2;"
const DELETE_MESSAGE_SQL = "DELETE FROM messages WHERE id = $1;"

type Message struct {
	Id string
	SendDate string
	Body string
	SenderId string
	ConversationId string
}

func InsertMessage(body, senderId, conversationId string) (Message, error) {
	var message Message

	err := DatabaseConnection.QueryRow(INSERT_MESSAGE_SQL, body, senderId, conversationId).Scan(&message.Id, &message.SendDate)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func SelectMessagesByConversation(conversationId string) ([]Message, error) {
	rows, err := DatabaseConnection.Query(SELECT_MESSAGES_BY_CONVERSATION_SQL, conversationId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	
	var messages []Message
	var message Message
	for rows.Next() {
		err = rows.Scan(&message.Id, &message.SendDate, &message.Body, &message.SenderId, &message.ConversationId)
		if err != nil {
			return nil, err
		}
		
		messages = append(messages, message)
	}

	return messages, nil
}

func (m *Message) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_MESSAGE_SQL, m.Body, m.Id)
	if err != nil {
		return err
	}
	
	return nil
}

func (m *Message) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_MESSAGE_SQL, m.Id)
	if err != nil {
		return err
	}

	return nil
}

