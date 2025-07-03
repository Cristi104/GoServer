package repository

const INSERT_MESSAGE_SQL = "INSERT INTO messages(body, sender_id, conversation_id) VALUES ($1, $2, $3) RETURNING id, send_date;"
const SELECT_MESSAGES_BY_CONVERSATION_SQL = "SELECT * FROM messages WHERE conversation_id = $1;"
const UPDATE_MESSAGE_SQL = "UPDATE messages SET body = $1 WHERE id = $2;"
const DELETE_MESSAGE_SQL = "DELETE FROM messages WHERE id = $1;"

type Message struct {
	id string
	sendDate string
	body string
	senderId string
	conversationId string
}

func InsertMessage(body, senderId, conversationId string) (Message, error) {
	var message Message

	err := DatabaseConnection.QueryRow(INSERT_MESSAGE_SQL, body, senderId, conversationId).Scan(&message.id, &message.sendDate)
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
		err = rows.Scan(&message.id, &message.sendDate, &message.body, &message.senderId, &message.conversationId)
		if err != nil {
			return nil, err
		}
		
		messages = append(messages, message)
	}

	return messages, nil
}

func (m *Message) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_MESSAGE_SQL, m.body, m.id)
	if err != nil {
		return err
	}
	
	return nil
}

func (m *Message) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_MESSAGE_SQL, m.id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) Id() string {
	return m.id
}

func (m *Message) SendDate() string {
	return m.sendDate
}

func (m *Message) Body() string {
	return m.body
}

func (m *Message) SenderId() string {
	return m.senderId
}

func (m *Message) ConversationId() string {
	return m.conversationId
}

func (m *Message) SetBody(body string)  {
	m.body = body
}

