# Queuetify
App to create a shared queue and a vote feature with the Spotify API.

## Installation
1) Install All Dependencies for the api. Make sure that you have golang install in your machine.

`cd queuetify/rest-api`

`go get .`

2) Copy the .env-example with the name .env and fill it. For the client keys [here](https://developer.spotify.com/documentation/general/guides/authorization/app-settings/).

`cp .env-example .env`

```
ADDRESS=127.0.0.1:8080 # rest api address
CORS=http://127.0.0.1:8081 # front-end address
CLIENT_ID=ahdfh-example-3jhjhssa # client id of the spotify api
CLIENT_SECRET=ajdsf-example-fj2lfhquh2iu3rh1jhds # client secret of the spotify api/authorization/
REDIRECT_URI=http://127.0.0.1:8081/app.html # redirect url of the front-end when create a room
CODE_SIZE=6 # size of the code to join the room
```
3) Install a static HTTP server for example: [http-server](https://www.npmjs.com/package/http-server)
 
4) Change the const in the /web-app/setup.js
```
const hostAddress = "120.0.0.1"
const hostBackEndPort = "8080"
const hostFrontEndPort = "8081"
```

## Usage
1) Run the rest-api `cd queuetify/rest-api` `go run *.go` or compile it `go build -o rest-api *.go` and `./rest-api
2) Launch the front-end example: `cd queuetify/web-app` `npx http-server .`
