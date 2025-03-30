window.onload = function(){
    let sidebar = document.getElementById("sidebar");
    fetch("/data/conversations", {
        method: "GET",
    })
    .then(response => response.json()) 
    .then(conversations => {
        let main = document.getElementById("main");
        conversations.forEach(conversation => {

            console.log(conversation)

            let button = document.createElement("button")
            button.classList += "conversation"
            button.id = conversation.Id
            button.innerHTML = conversation.Name

            button.addEventListener("click", function(event) {
                event.stopPropagation();

                main.innerHTML = ""

                fetch("/data/messages?id=" + button.id, {
                    method: "GET"
                })
                .then(response => response.json())
                .then(messages => {
                    console.log(messages)
                    messages.forEach(Message => {
                        let message = document.createElement("div")
                        message.classList += "message"

                        let name = document.createElement("h3")
                        name.innerHTML = conversation.Members.find((e) => e.Id == Message.SenderId).Username
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

                        main.append(message)
                    })
                })
            })
            sidebar.append(button)
            
        });
    })
    .catch(error => console.error('Error:', error));
}