window.onload = function(){
    fetch("/data/conversations", {
        method: "GET",
    })
    .then(response => response.json()) 
    .then(conversations => {
        conversations.forEach(conversation => {
            addConversation(conversation)
        });
    })
    .catch(error => console.error('Error:', error));
}

function addConversation(conversation){
    let main = document.getElementById("main");
    let sidebar = document.getElementById("sidebar");

    let button = document.createElement("button")
    button.classList += "conversation"
    button.id = conversation.Id
    button.innerHTML = conversation.Name

    button.addEventListener("click", function(event) {
        event.stopPropagation();

        main.innerHTML = ""
        let data = new FormData()
        data.append("ConversationId", conversation.Id)

        fetch("/data/messages/load", {
            method: "POST",
            body: data,
        })
        .then(response => response.json())
        .then(messages => {
            createMessageForm(conversation)
            messages.forEach(Message => {
                addMessage(conversation.Members.find((e) => e.Id == Message.SenderId).Username, Message)
            })
        })
    })
    sidebar.appendChild(button)
}

function addMessage(username, Message){
    let main = document.getElementById("main")
    let message = document.createElement("div")
    message.classList += "message"

    let name = document.createElement("h3")
    name.innerHTML = username //conversation.Members.find((e) => e.Id == Message.SenderId).Username
    name.classList += "username"
    message.append(name)

    let messageBody = document.createElement("p")
    messageBody.innerHTML = Message.Body
    messageBody.classList += "messageBody"
    message.append(messageBody)

    let date = document.createElement("p")
    date.innerHTML = Message.SendDate
    date.classList += "date"
    message.append(date)

    main.insertBefore(message, main.lastChild)
}

function createMessageForm(conversation){
    let main = document.getElementById("main")
    let messageForm = document.createElement("form")
    messageForm.id = "messageForm"
    messageForm.action = "/data/messages"
    messageForm.method = "POST"

    let messageInput = document.createElement("input")
    messageInput.type = "text"
    messageInput.id = "message"
    messageInput.name = "message"
    messageInput.required = true
    messageForm.appendChild(messageInput)

    let messageSubmit = document.createElement("input")
    messageSubmit.type = "submit"
    messageForm.appendChild(messageSubmit)

    messageForm.addEventListener("submit", function(event){
        event.preventDefault()
        let data = new FormData(messageForm)
        data.append("ConversationId", conversation.Id)
        fetch("/data/messages/add", {
            method: "POST",
            body: data,
        })
        .then(response => response.json())
        .then(messages => {
            messages.forEach(Message => {
                addMessage(conversation.Members.find((e) => e.Id == Message.SenderId).Username, Message)
            })
        })
        .catch(error => console.error('Error:', error));
    })

    main.append(messageForm)

}