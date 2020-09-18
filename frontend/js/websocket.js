const info = document.querySelector('.shadowCont')
const ctx = document.getElementById("matchfield").getContext("2d")
const fieldSize = 50;
const canvasSize = 1000;
let socket = new WebSocket("ws://localhost:2100/ws-test/")
let playerChar = new Image();
let ticker;
let isBombLegal = true;
let keyPresses = {};
let userId;
let currentUser;

let wallImg = new Image();
let wallImg2 = new Image();
let grassImg = new Image();
let bombImg = new Image();
let itemBoostImg = new Image()
let itemSlowImg = new Image()
let itemGhostImg = new Image()
let playerGhostImg = new Image()

const nameLabel = document.createElement("p");
const posXLabel = document.createElement("p");
const posYLabel = document.createElement("p");

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


    if (ctx !== null && incomingPackage !== null){
        ctx.clearRect(0, 0, canvasSize, canvasSize);
        background(grassImg, incomingPackage.GameMap);
        searchForUser(incomingPackage.Players)
        updateUserInfo()


        for(let i = 0; i < incomingPackage.Players.length; i++){
            ctx.fillText(incomingPackage.Players[i].Name,incomingPackage.Players[i].PositionX + 15,incomingPackage.Players[i].PositionY - 5, 100);

            if (incomingPackage.Players[i] != null && !incomingPackage.Players[i].GhostActive){
                drawPlayerPosClient();
                ctx.drawImage(playerChar, incomingPackage.Players[i].PositionX, incomingPackage.Players[i].PositionY, fieldSize, fieldSize);
            }
        }
        drawImageFromEnum(itemBoostImg, gamemap, 6);
        drawImageFromEnum(itemSlowImg, gamemap, 7);
        drawImageFromEnum(itemGhostImg, gamemap, 8);
        drawImageFromEnum(wallImg, incomingPackage.GameMap, 3);
        drawImageFromEnum(wallImg2, incomingPackage.GameMap, 2);
        drawImageFromEnum(bombImg, gamemap, 1);


        for(let i = 0; i < incomingPackage.Players.length; i++){
            ctx.fillText(incomingPackage.Players[i].Name,incomingPackage.Players[i].PositionX + 15,incomingPackage.Players[i].PositionY - 5, 100);
            if(incomingPackage.Players[i] != null && incomingPackage.Players[i].GhostActive){
                drawPlayerPosClient();
                ctx.globalAlpha = 0.5
                if (incomingPackage.Players[i] != null && incomingPackage.Players[i].IsAlive){
                    ctx.drawImage(playerChar, incomingPackage.Players[i].PositionX, incomingPackage.Players[i].PositionY, fieldSize, fieldSize);
                }
                else {
                    ctx.drawImage(playerGhostImg, incomingPackage.Players[i].PositionX, incomingPackage.Players[i].PositionY, fieldSize, fieldSize);
                }
                ctx.globalAlpha = 1
            }
        }
    }

}

function initGame(){
    playerChar.src = "media/player.png"
    wallImg.src ="media/wallBreak.png"
    wallImg2.src ="media/wallBreak2.png"
    grassImg.src = "media/grass.png"
    bombImg.src = "media/bomb.png"
    playerChar.src = "media/cutie2.png"
    itemBoostImg.src = "media/speeditem.png"
    itemSlowImg.src = "media/slowitem.png"
    itemGhostImg.src = "media/ghostitem.png"
    playerGhostImg.src = "media/ghostPlayer.png"
    info.append(nameLabel);
    info.append(posXLabel);
    info.append(posYLabel);
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
function drawImageFromEnum (img, map, type){
    for (i = 0; i < map.length; i++){
        for (j = 0; j < map[i].length; j++) {
            for (k = 0; k < map[i][j].length; k++){
                if (map[i][j][k] === type) {
                    ctx.drawImage(img, i *fieldSize, j * fieldSize, fieldSize, fieldSize);
                }
            }
        }
    }
}
function background (img, map){
    for (i = 0; i < map.length; i++){
        for (j = 0; j < map[i].length; j++) {
                         ctx.drawImage(img, i *fieldSize, j * fieldSize, fieldSize, fieldSize);
        }
    }
}

function updateUserInfo() {
    if(currentUser == null || nameLabel == null || posXLabel == null || posYLabel == null){return}
    nameLabel.innerHTML = "Name: " + currentUser.Name;
    posXLabel.innerHTML = "Position X: " + currentUser.PositionX;
    posYLabel.innerHTML = "Position Y: " + currentUser.PositionY;
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

function searchForUser(playerArr){
    if (userId === "" || playerArr == null){
        return
    }
    for(let i = 0; i < playerArr.length; i++) {
        if (playerArr[i].UserID == userId) {
            currentUser = playerArr[i]
            return
        }
    }
}

function drawPlayerPosClient() {
    if(currentUser == null){return}
    ctx.fillStyle = "rgba(0, 0, 0, 0.4)";
    const x = Math.floor((currentUser.PositionX + fieldSize/2)/fieldSize) * fieldSize
    const y = Math.floor((currentUser.PositionY + fieldSize/2)/fieldSize) * fieldSize
    ctx.fillRect(x,y , fieldSize , fieldSize)
}


document.addEventListener('keydown', keyDownListener, false);
document.addEventListener('keydown', spaceKeyDownListener, false);
document.addEventListener('keyup', keyUpListener, false);

function spaceKeyDownListener(event) {
    if(isBombLegal && event.key === " "){
        keyPresses[event.key] = true
        isBombLegal = false
        setTimeout(()=>{
            isBombLegal = true
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
    var canvas = document.getElementById("matchfield");
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

