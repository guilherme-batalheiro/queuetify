function navigateToPage(url) {
    const fullUrl = window.location.origin + url;
    
    window.location.href = fullUrl;
}

let el = document.getElementById("loginButton") 
if (el) {
    el.addEventListener("click", function() {
        navigateToPage('/userSpotifyLogin');
    });
}

el = document.getElementById("createRoomButton") 
if (el) {
    el.addEventListener("click", function() {
        navigateToPage('/room/createRoom');
    });
}


