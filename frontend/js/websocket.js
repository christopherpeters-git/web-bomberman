let testContainer = document.getElementById("test");
const ctx = document.getElementById("matchfield").getContext("2d")
let socket = new WebSocket("ws://localhost:2100/ws-test/")
let playerChar = new Image();
let ticker;
let keyPresses = {};
playerChar.src = "media/player1.png"
console.log("Attempting Websocket connection")

console.log(testContainer);

socket.onopen = () => {
    ticker = setInterval(function(){socket.send(JSON.stringify(keyPresses))}, 5);
    console.log("Connected")
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}

socket.onmessage = (ev) => {
    const incomingPackage = JSON.parse(ev.data)

    if (ctx !== null && incomingPackage !== null){
        ctx.clearRect(0, 0, 500, 500);
        drawGrid(500, 500, "matchfield");

        for(let i = 0; i < incomingPackage.Players.length; i++){
            ctx.fillText(incomingPackage.Players[i].Name,incomingPackage.Players[i].PositionX + 15,incomingPackage.Players[i].PositionY - 5, 100);
            ctx.drawImage(playerChar, incomingPackage.Players[i].PositionX, incomingPackage.Players[i].PositionY, 50, 50);
        }
    }
    testContainer.innerHTML =
        "UserID: " + incomingPackage.Players[0].UserID + "<br>"
        + "Username: " + incomingPackage.Players[0].Name + "<br>"
        + "X-Postion: " + incomingPackage.Players[0].PositionX + "<br>"
        + "Y-Postion: " + incomingPackage.Players[0].PositionY + "<br>"
        + "Spieler lebt: " + incomingPackage.Players[0].IsAlive + "<br>"
        + "Bombenradius: " + incomingPackage.Players[0].BombRadius + "<br>" + "<br>";


    drawElement("#ae1111",incomingPackage.GameMap, 2 )
    drawElement("#000000",incomingPackage.GameMap, 1 )
}

function drawElement (color, map, type){
    for (i = 0; i < map.length; i++){
        for (j = 0; j < map[i].length; j++) {
            for (k = 0; k < map[i][j].length; k++){
                if (map[i][j][k] === type) {
                    ctx.fillStyle = color;
                    ctx.fillRect(i * 50, j * 50, 50 , 50)
                }
            }
        }
    }
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

