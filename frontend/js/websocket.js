// let testContainer = document.getElementById("test");
const ctx = document.getElementById("matchfield").getContext("2d")
const fieldSize = 50
let socket = new WebSocket("ws://localhost:2100/ws-test/")
let playerChar = new Image();
let ticker;
let isBombLegal = true;
let keyPresses = {};
let userId = "";
playerChar.src = "media/player1.png"
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
            console.log(userId)
        })
        .catch((reason => {
            console.log(reason)
        }))
    console.log("Connected")
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}

socket.onmessage = (ev) => {
    const incomingPackage = JSON.parse(ev.data);

    if (ctx !== null && incomingPackage !== null){
        ctx.clearRect(0, 0, 500, 500);
        drawGrid(500, 500, "matchfield");

        for(let i = 0; i < incomingPackage.Players.length; i++){
            ctx.fillText(incomingPackage.Players[i].Name,incomingPackage.Players[i].PositionX + 15,incomingPackage.Players[i].PositionY - 5, 100);
            ctx.drawImage(playerChar, incomingPackage.Players[i].PositionX, incomingPackage.Players[i].PositionY, 50, 50);
        }
    }

    drawElement("#ae1111",incomingPackage.GameMap, 3 )
    drawElement("#60f542",incomingPackage.GameMap, 2 )
    drawElement("#000000",incomingPackage.GameMap, 1 )
    drawPlayerPosClient(incomingPackage.Players)
    // drawPlayersPos(incomingPackage.TestPlayer)
}

function drawElement (color, map, type){
    for (let i = 0; i < map.length; i++){
        for (let j = 0; j < map[i].length; j++) {
            for (let k = 0; k < map[i][j].length; k++){
                if (map[i][j][k] === type) {
                    ctx.fillStyle = color;
                    ctx.fillRect(i * fieldSize, j * fieldSize, fieldSize ,fieldSize)
                }
            }
        }
    }
}
function drawPlayersPos(playerArr) {
    for (let i = 0; i < playerArr.length; i++){
        for (let j = 0; j < playerArr[i].length; j++) {
                if (playerArr[i][j] === 1) {
                    ctx.fillStyle = "rgba(0, 0, 0, 0.4)";
                    ctx.fillRect(i * fieldSize, j * fieldSize, fieldSize , fieldSize)
                }
            }
        }
}

function drawPlayerPosClient(playerArr) {
    if (userId === "" || playerArr == null){
        return
    }
    for(let i = 0; i < playerArr.length; i++){
        if(playerArr[i].UserID == userId){
            ctx.fillStyle = "rgba(0, 0, 0, 0.4)";
            const x = Math.floor((playerArr[i].PositionX + fieldSize/2)/fieldSize) * fieldSize
            const y = Math.floor((playerArr[i].PositionY + fieldSize/2)/fieldSize) * fieldSize
            ctx.fillRect(x,y , fieldSize , fieldSize)
            break
        }
    }
}


document.addEventListener('keydown', keyDownListener, false);
document.addEventListener('keydown', spaceKeyDownListener, false);
document.addEventListener('keyup', keyUpListener, false);

function spaceKeyDownListener(event) {
    if(isBombLegal && event.key === " "){
        keyPresses[event.key] = true
        isBombLegal = false
        console.log("lol1")
        setTimeout(()=>{
            isBombLegal = true
            console.log("lol2")
        },1000)
    }
}


function keyDownListener(event) {
    if(event.key !== " "){keyPresses[event.key] = true;}
}

function keyUpListener(event) {
    keyPresses[event.key] = false;
}




var drawGrid = function(w, h, id) {
    var canvas = document.getElementById(id);
    var ctx = canvas.getContext('2d');
    ctx.canvas.width = w;
    ctx.canvas.height = h;


    for (x = 0; x <= w; x += fieldSize) {
        ctx.moveTo(x, 0);
        ctx.lineTo(x, h);
        for (y = 0; y <= h; y += fieldSize) {
            ctx.moveTo(0, y);
            ctx.lineTo(w, y);
        }
    }
    ctx.stroke();
};

