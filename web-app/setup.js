const hostAddress = "127.0.0.1"
const hostBackEndPort = "8080"
const hostFrontEndPort = "8081"
const clientId = ""

const url = "http://accounts.spotify.com/authorize?response_type=code&client_id=" + clientId + "&scope=user-read-email%20user-modify-playback-state%20user-read-playback-state%20user-read-currently-playing&redirect_uri=http%3A%2F%2F" + hostAddress + "%3A" + hostFrontEndPort + "%2Fapp.html"

const urlSearchParams = new URLSearchParams(window.location.search)
const params = Object.fromEntries(urlSearchParams.entries())
const token = params.code

function createRoom() {
    // Creates room by log in with spotify, because just the spotify users can create rooms.
    window.open(url, "_self")
}

function startAuth(autherizationCode) {
    // This funtion create room and log in into the app.
    fetch('http://' + hostAddress + ':' + hostBackEndPort + '/auth?code=' + autherizationCode, {
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
            
        fetch('http://' + hostAddress + ':' + hostBackEndPort + '/create_room?user_id=' + data.user_id, {
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
            window.open("http://" + hostAddress + ":" + hostFrontEndPort + "/app.html", "_self")
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
    fetch('http://' + hostAddress + ':' + hostBackEndPort + '/delete_room?user_id=' + localStorage.getItem("user_id"), {
        method: 'post',
    })
    .then(response => {
        
        localStorage.clear() 
        window.open("http://" + hostAddress + ":" + hostFrontEndPort + "/app.html", "_self")
    })
}

function exitRoom() {
    // Exit room if the request failed close the room anywais
    fetch('http://' + hostAddress + ':' + hostBackEndPort + '/exit_room?room_code=' + localStorage.getItem("room_code"), {
        method: 'post',
    })
    .then(response => {
        localStorage.clear() 
        window.open("http://" + hostAddress + ":" + hostFrontEndPort + "/app.html", "_self")
    })
}

function addMusicToQueue() {
    // Add music to queue
    let songName = document.getElementById('songNameInput').value
    fetch('http://' + hostAddress + ':' + hostBackEndPort + '/add_to_queue?room_code=' + localStorage.getItem('room_code') + '&song_name=' + songName, {
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
        fetch('http://' + hostAddress + ':' + hostBackEndPort + '/join_room?room_code=' + room_code, {
            method: 'post',})
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            localStorage.setItem('room_code', room_code);
            window.open("http://" + hostAddress + ":" + hostFrontEndPort + "/app.html", "_self")
        })
        .catch(function(error) {
            alert("Failed to join room!")
            console.log(error);
        });
    }
}

function updateSong(room_code, loop) {
    fetch('http://' + hostAddress + ':' + hostBackEndPort + '/current_song?room_code=' + room_code, {
        method: 'get',})
    .then(response => {
        if (!response.ok) {
            setTimeout(updateSong, 5000, room_code, true)
            throw Error(response.statusText);
        } else if (response.status == 204) {
            document.getElementById("song_name").innerHTML = localStorage.getItem('song_name')
            document.getElementById("song_artits").innerHTML = localStorage.getItem('song_artits')
            setTimeout(updateSong, 50000, room_code, true)
            throw Error(response.statusText);
        }

        return response.json()
    })
    .then(data => {
        console.log(data.song_name)
        if (data.song_name != localStorage.getItem("song_name")) {
            localStorage.setItem("song_name", data.song_name)
            localStorage.setItem("song_artists", data.song_artists)
            localStorage.setItem("can_vote", "true")
        }
        if(loop) {
            let timer = data.duration_ms - data.progress_ms + 500
            setTimeout(updateSong, timer, room_code, true)
        }
        document.getElementById("song_name").innerHTML = localStorage.getItem('song_name')
        document.getElementById("song_artists").innerHTML = localStorage.getItem('song_artists')
    })
    .catch(function(error) {
        console.log(error);
    });
}

function voteToSkipSong() {
	console.log(localStorage.getItem("can_vote"))
    if(localStorage.getItem("can_vote") == "true") {
	    if (progress_voting) {
            return;
        }

        progress_voting = true;

        fetch('http://' + hostAddress + ':' + hostBackEndPort + '/vote_skip_song?room_code=' + localStorage.getItem("room_code"), {
            method: 'get',})
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
			localStorage.setItem("can_vote", "false")
            return response.json()
        })
        .then(data => {
            if(!data.music_skipped) 
                alert(data.missing_votes + " people left to vote!")
            else
                alert("Music skipped sucefully!")
            updateSong(localStorage.getItem("room_code"), false)
        })
        .catch(function(error) {
            alert("Vote failed!")
			localStorage.setItem("can_vote", "true")
            console.log(error);
        })
        .finally(() => {
            progress_voting = false;
        });

    } else {
		alert("You already voted if a new song is playing refresh the page.")
	}
}

if (token != null) {
    startAuth(token)
}

let room_code = localStorage.getItem("room_code") 
let progress_voting = false;
if (room_code != null) {
    updateSong(room_code, true)
}


