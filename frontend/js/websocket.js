let testContainer = document.getElementById("test");
const ctx = document.getElementById("matchfield").getContext("2d")
let socket = new WebSocket("ws://localhost:2100/ws-test/")
let playerChar = new Image();
let keyPresses = {};
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
        drawGrid(500, 500, "matchfield");

        for(let i = 0; i < users.length; i++){
            ctx.fillText(users[i].Name,users[i].PositionX + 15,users[i].PositionY - 5, 100);
            //ctx.fillRect(users[i].PositionX, users[i].PositionY, 50, 50);
            ctx.drawImage(playerChar, users[i].PositionX, users[i].PositionY, 50, 50);
        }
    }
    testContainer.innerHTML = ev.data;
}

document.addEventListener('keydown', keyDownListener, false);
document.addEventListener('keyup', keyUpListener, false);

function keyDownListener(event) {
    keyPresses[event.key] = true;
}

function keyUpListener(event) {
    keyPresses[event.key] = false;
}

var drawGrid = function(w, h, id) {
    var canvas = document.getElementById(id);
    var ctx = canvas.getContext('2d');
    ctx.canvas.width = w;
    ctx.canvas.height = h;


    for (x = 0; x <= w; x += 50) {
        ctx.moveTo(x, 0);
        ctx.lineTo(x, h);
        for (y = 0; y <= h; y += 50) {
            ctx.moveTo(0, y);
            ctx.lineTo(w, y);
        }
    }
    ctx.stroke();
};

const ticker = setInterval(function(){socket.send(JSON.stringify(keyPresses))}, 5);
