const url = "http://accounts.spotify.com/authorize?response_type=code&client_id=f4a290e1e7a648458eb0cc169279c4ac&scope=user-read-email%20user-modify-playback-state%20user-read-playback-state%20user-read-currently-playing&redirect_uri=http%3A%2F%2F127.0.0.1%3A8081%2Fapp.html"

const urlSearchParams = new URLSearchParams(window.location.search)
const params = Object.fromEntries(urlSearchParams.entries())
const token = params.code

function createRoom() {
    // Creates room by log in with spotify, because just the spotify users can create rooms.
    window.open(url, "_self")
}

function startAuth(autherizationCode) {
    // This funtion create room and log in into the app.
    fetch('http://127.0.0.1:8080/auth?code=' + autherizationCode, {
        method: 'get',
    })
    .then(response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }

        return response.json()
    })
    .then(data => {
        localStorage.setItem('user_id', data.user_id);
        localStorage.setItem('display_name', data.display_name);
            
        fetch('http://127.0.0.1:8080/create_room?user_id=' + data.user_id, {
            method: 'get',
        })
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }

            return response.json()
        })
        .then(data => { 
            localStorage.setItem('room_code', data.room_code);
            window.open("http://127.0.0.1:8081/app.html", "_self")
        })
        .catch(function(error) {
            console.log(error);
            alert("Authorization failed!")
        });

    })
    .catch(function(error) {
        console.log(error);
        alert("Authorization failed!")
    });
}

function closeRoom() {
    // Close room if the request failed close the room anywais
    fetch('http://127.0.0.1:8080/delete_room?user_id=' + localStorage.getItem("user_id"), {
        method: 'post',
    })
    .then(response => {
        
        localStorage.clear() 
        window.open("http://127.0.0.1:8081/app.html", "_self")
    })
}

function exitRoom() {
    // Exit room if the request failed close the room anywais
    fetch('http://127.0.0.1:8080/exit_room?room_code=' + localStorage.getItem("room_code"), {
        method: 'post',
    })
    .then(response => {
        localStorage.clear() 
        window.open("http://127.0.0.1:8081/app.html", "_self")
    })
}

function addMusicToQueue() {
    // Add music to queue
    let songName = document.getElementById('songNameInput').value
    fetch('http://127.0.0.1:8080/add_to_queue?room_code=' + localStorage.getItem('room_code') + '&song_name=' + songName, {
        method: 'get',})
    .then(response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }
        return response.json()
    })
    .then(data => {
        alert("Music: " +  data.song_name + "\nBy artirts: " +  data.song_artists.toString() + "\nSusseful added to queue!") 
    })
    .catch(function(error) {
        alert("Failed to add to queue!")
        console.log(error);
    });
}

function joinRoom() {
    let room_code = prompt("Please a room code:","******");
    if (room_code != null && room_code.length == 6) {
        fetch('http://127.0.0.1:8080/join_room?room_code=' + room_code, {
            method: 'post',})
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            localStorage.setItem('room_code', room_code);
            window.open("http://127.0.0.1:8081/app.html", "_self")
        })
        .catch(function(error) {
            alert("Failed to join room!")
            console.log(error);
        });
    }
}

function updateSong(room_code, loop) {
    fetch('http://127.0.0.1:8080/current_song?room_code=' + room_code, {
        method: 'get',})
    .then(response => {
        if (!response.ok) {
            setTimeout(updateSong, 5000, room_code, true)
            throw Error(response.statusText);
        } else if (response.status == 204) {
            console.log(data.song_name)
            setTimeout(updateSong, 50000, room_code, true)
            throw Error(response.statusText);
        }

        return response.json()
    })
    .then(data => {
        console.log(data.song_name)
        if (data.song_name != localStorage.getItem("song_name")) {
            localStorage.setItem("song_name", data.song_name)
            localStorage.setItem("can_vote", true)
        }
        if(loop) {
            let timer = data.duration_ms - data.progress_ms + 500
            setTimeout(updateSong, timer, room_code, true)
        }
    })
    .catch(function(error) {
        console.log(error);
    });
}

function voteToSkipSong() {
    if(localStorage.getItem("can_vote")) {
        fetch('http://127.0.0.1:8080/vote_skip_song?room_code=' + localStorage.getItem("room_code"), {
            method: 'get',})
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            return response.json()
        })
        .then(data => {
            console.log(data)
            if(!data.music_skipped) 
                alert(data.missing_votes + " people left to vote!")
            else
                alert("Music skipped sucefully!")
            updateSong(localStorage.getItem("room_code"), false)
        })
        .catch(function(error) {
            alert("Vote failed!")
            console.log(error);
        });

    }
}

if (token != null) {
    startAuth(token)
}

let room_code = localStorage.getItem("room_code") 
if (room_code != null) {
    updateSong(room_code, true)
}


