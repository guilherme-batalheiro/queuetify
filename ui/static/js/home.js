function navigateToPage(url) {
    const fullUrl = window.location.origin + url;
    
    window.location.href = fullUrl;
}

document.getElementById("loginButton").addEventListener("click", function() {
    navigateToPage('/userSpotifyLogin/');
});
