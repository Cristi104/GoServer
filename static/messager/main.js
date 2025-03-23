window.onload = function(){
    let sidebar = document.getElementById("sidebar")
    fetch("/data/friends", {
        method: "GET",
    })
    .then(response => response.json()) 
    .then(data => {
        sidebar.innerHTML += "<p>" + data + "</p>";
    })
    .catch(error => console.error('Error:', error));
}