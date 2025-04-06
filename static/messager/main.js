window.onload = function(){
    loadAccountData()
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

function loadAccountData(){
    fetch("/data/account", {
        method: "GET",
    })
    .then(response => response.json())
    .then(accountData => {
        localStorage.setItem("account", JSON.stringify(accountData))
        console.log(accountData)
    })
    .catch(error => console.error('Error:', error));
}

function addConversation(conversation){
    let main = document.getElementById("messages");
    let sidebar = document.getElementById("sidebar");

    let button = document.createElement("button")
    button.classList.add("conversation")
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
    let main = document.getElementById("messages")
    let message = document.createElement("div")
    let account = JSON.parse(localStorage.getItem("account"))
    console.log(account.Id)
    console.log(Message.SenderId)
    if(Message.SenderId == account.Id)
        message.classList.add("userMessage")
    message.classList.add("message")

    let name = document.createElement("h3")
    name.innerHTML = username
    name.classList.add("username")
    message.append(name)

    let messageBody = document.createElement("p")
    messageBody.innerHTML = Message.Body
    messageBody.classList.add("messageBod")
    message.append(messageBody)

    let date = document.createElement("p")
    date.innerHTML = Message.SendDate
    date.classList.add("date")
    message.append(date)

    main.append(message)
}

function createMessageForm(conversation){
    // let main = document.getElementById("content")
    // let messageForm = document.createElement("form")
    // messageForm.id = "messageForm"
    // messageForm.action = "/data/messages"
    // messageForm.method = "POST"

    // let messageInput = document.createElement("input")
    // messageInput.type = "text"
    // messageInput.id = "message"
    // messageInput.name = "message"
    // messageInput.autocomplete = "off"
    // messageInput.required = true
    // messageForm.appendChild(messageInput)

    // let messageSubmit = document.createElement("input")
    // messageSubmit.type = "submit"
    // messageSubmit.value = "send"
    // messageForm.appendChild(messageSubmit)

    let messageForm = document.getElementById("messageForm")
    messageForm.addEventListener("submit", function(event){
        event.preventDefault()
        let data = new FormData(messageForm)
        data.append("ConversationId", conversation.Id)
        let messageInput = document.getElementById("message")
        messageInput.value = ""
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

    // main.append(messageForm)
}