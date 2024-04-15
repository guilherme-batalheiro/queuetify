function navigateToPage(url) {
    const fullUrl = window.location.origin + url;
    
    window.location.href = fullUrl;
}

let el = document.getElementById("deleteRoomButton") 

if (el) {
    el.addEventListener("click", function() {
        navigateToPage(window.location.pathname+'/deleteRoom');
    });
}

el = document.getElementById("addSongToQueue")
if (el) {
    el.addEventListener("click", function() { addToQueue(); })
}

function addToQueue() {
    const data = new URLSearchParams();
    const song_query = document.getElementById("musicInput").value;
    data.append('song_query', song_query);
    
    fetch(window.location.pathname+'/addToQueue', {
        method: "post",
        body: data,
    })
    .then(response => console.log(response))
    .catch(error => {
        console.error("Error:", error);
    });
}
