function updateSong(room_code) {
    setTimeout(() => { }, 5000)
    postMessage(room_code)
}

self.addEventListener("message", function(e) {
    room_code = e.data
    updateSong(room_code)
}, false);
