let testContainer = document.getElementById("test");
const ctx = document.getElementById("matchfield").getContext("2d")
let socket = new WebSocket("ws://localhost:2100/ws-test/")
let playerChar = new Image();
playerChar.src = "media/player1.png"
console.log("Attempting Websocket connection")

console.log(testContainer);

socket.onopen = () => {
    console.log("Connected")
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}

socket.onmessage = (ev) => {
    const users = JSON.parse(ev.data)

    if (ctx !== null && users !== null){
        ctx.clearRect(0, 0, 500, 500);

        for(let i = 0; i < users.length; i++){
            ctx.fillText(users[i].Name,users[i].PositionX + 15,users[i].PositionY - 5, 100);
            //ctx.fillRect(users[i].PositionX, users[i].PositionY, 50, 50);
            ctx.drawImage(playerChar, users[i].PositionX, users[i].PositionY, 50, 50);
        }
    }
    testContainer.innerHTML = ev.data;
}

document.addEventListener( 'keydown', handleKeyPress, false );

function handleKeyPress(event){
    let keyCode;
    if (event.key !== undefined) {
        keyCode = event.key;
    } else if (event.keyIdentifier !== undefined) {
        keyCode = event.keyIdentifier;
    } else if (event.keyCode !== undefined) {
        keyCode = event.keyCode;
    }
    socket.send(keyCode);
}