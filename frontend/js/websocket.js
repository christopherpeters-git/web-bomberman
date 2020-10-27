let socket = new WebSocket("ws://localhost:2100/ws-test/")
let ticker;
let userId;
let exec = false;
let intervall;
initGame();

console.log("Attempting Websocket connection")
socket.onopen = () => {
    ticker = setInterval(function(){
        socket.send(JSON.stringify(keyPresses))
        keyPresses[" "] = false;
    }, 10);
    fetch("/fetchUserId")
        .then(response => {
            if(response.status === 200){
               return response.text()
            }
        })
        .then(text => {
            userId = text
            console.log("User ID: " + userId)
        })
        .catch((reason => {
            console.log(reason)
        }))
    console.log("Connected")
    frameCounter = 0;
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}

socket.onmessage = (ev) => {
    const incomingPackage = JSON.parse(ev.data);
    const gamemap = incomingPackage.GameMap;
    const players = incomingPackage.Players;
    let sessionRunning = incomingPackage.SessionRunning;

    updateUI(sessionRunning);

    if (ctx !== null && incomingPackage !== null){
       gameLoop(gamemap, players);
    }
}











