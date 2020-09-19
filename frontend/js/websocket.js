const info = document.querySelector('.shadowCont')
const ctx = document.getElementById("matchfield").getContext("2d")
const fieldSize = 50;
const canvasSize = 1000;
let socket = new WebSocket("ws://localhost:2100/ws-test/")

let playerChar = new Image();
let playerChar2 = new Image();
let playerChar3 = new Image();

let playerCharUp = new Image();
let playerCharLeft = new Image();
let playerCharRight = new Image();
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

let counter = 0;
let imgIndex = 0;
const frameLimit = 2;

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
        counter++;
        ctx.clearRect(0, 0, canvasSize, canvasSize);
        background(grassImg, incomingPackage.GameMap);
        searchForUser(incomingPackage.Players)
        updateUserInfo()


        for(let i = 0; i < incomingPackage.Players.length; i++){
            ctx.fillText(incomingPackage.Players[i].Name,incomingPackage.Players[i].PositionX + 15,incomingPackage.Players[i].PositionY - 5, 100);
            if (incomingPackage.Players[i] != null && !incomingPackage.Players[i].GhostActive){
                drawPlayerPosClient();
                drawPlayerChar(incomingPackage.Players[i], counter)
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
                drawPlayerChar(incomingPackage.Players[i], counter)
                ctx.globalAlpha = 1
            }
        }

        if (counter == frameLimit){
            counter = 0;
        }

    }
}

function drawPlayerChar (player, count) {

    playerImgDown = playerChar;

    if (count == frameLimit) {
        imgIndex++;
    }
    if (imgIndex == 4) {
        imgIndex = 0;
    }
    if (imgIndex == 0) {
        playerImgDown = playerChar;
    }
    if(imgIndex == 1) {
        playerImgDown = playerChar2;
    }
    if(imgIndex == 2) {
        playerImgDown = playerChar;
    }
    if(imgIndex == 3) {
        playerImgDown = playerChar3;
    }

    if (player.IsAlive){
        if (player.DirDown){
            ctx.drawImage(playerImgDown, player.PositionX, player.PositionY, fieldSize, fieldSize);
        } else if ( player.DirUp){
            ctx.drawImage(playerCharUp, player.PositionX, player.PositionY, fieldSize, fieldSize);
        } else if (player.DirLeft){
            ctx.drawImage(playerCharLeft, player.PositionX, player.PositionY, fieldSize, fieldSize);
        } else if (player.DirRight){
            ctx.drawImage(playerCharRight, player.PositionX, player.PositionY, fieldSize, fieldSize);
        }else {
            ctx.drawImage(playerChar, player.PositionX, player.PositionY, fieldSize, fieldSize);
        }
    }
    else {
        ctx.drawImage(playerGhostImg, player.PositionX, player.PositionY, fieldSize, fieldSize);
    }
}


function initGame(){
    wallImg.src ="media/wallBreak.png"
    wallImg2.src ="media/wallBreak2.png"
    grassImg.src = "media/grass.png"
    bombImg.src = "media/bomb.png"
    playerChar.src = "media/cutie0.png"
    playerCharUp.src  = "media/cutieBack.png"
    playerCharLeft.src  = "media/cutieLeft.png"
    playerCharRight.src  = "media/cutieRight.png"
    itemBoostImg.src = "media/speeditem.png"
    itemSlowImg.src = "media/slowitem.png"
    itemGhostImg.src = "media/ghostitem.png"
    playerGhostImg.src = "media/ghostPlayer.png"
    playerChar2.src = "media/cutie1.png"
    playerChar3.src = "media/cutie3.png"
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


