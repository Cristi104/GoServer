window.onload = function(){
    let sidebar = document.getElementById("sidebar");
    fetch("/data/friends", {
        method: "GET",
    })
    .then(response => response.json()) 
    .then(data => {
        let main = document.getElementById("main");
        data.forEach(element => {
            sidebar.innerHTML += "<button class=\"conversation\" id=\"" + element + "\">" + element + "</button>";
            let button = document.getElementById(element)
            button.addEventListener("click", function(event) {
                event.stopPropagation();
                fetch("/data/conversation?id=" + button.id, {
                    method: "GET"
                })
                .then(response => response.json())
                .then(data => {
                    main.innerHTML = ""
                    console.log(data)
                    data.forEach(dataElement => {
                        main.innerHTML += "<p>" + JSON.stringify(dataElement) + "</p>"
                    })
                })
            })
            
        });
    })
    .catch(error => console.error('Error:', error));

    // let conversations = document.querySelectorAll(".conversation");
    // let main = document.getElementById("main");
    // console.log(conversations)
    // conversations.forEach(element => {
    //     element.addEventListener("click", function(event) {
    //         event.stopPropagation();
    //         fetch("/data/conversation?id=" + element.id, {
    //             method: "GET"
    //         })
    //         .then(response => response.json())
    //         .then(data => {
    //             main.innerHTML = ""
    //             data.forEach(dataElement => {
    //                 main.innerHTML += "<p>" + dataElement + "</p>"
    //             })
    //         })
    //     })
    // })
}