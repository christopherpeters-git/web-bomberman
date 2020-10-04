//Fieldobject Images
let playerCharSprite = new Image();
let wallImg = new Image();
let wallImg2 = new Image();
let grassImg = new Image();
let bombImg = new Image();
let bomb2Img = new Image();
let bomb3Img = new Image();
let itemBoostImg = new Image();
let itemSlowImg = new Image();
let itemGhostImg = new Image();
let playerGhostImg = new Image();
let explosionImg = new Image();
let portalImg = new Image();
let poisonImg = new Image();

let frameCounter = 0;
let imgIndex = 0;


function initGame(){
    wallImg.src ="media/hardwall.png"
    wallImg2.src ="media/weakwall.png"
    grassImg.src = "media/metalfloor.png"
    bombImg.src = "media/bomb3.png"
    bomb2Img.src ="media/bombstate2.png"
    bomb3Img.src ="media/bombstate3.png"
    playerCharSprite.src ="media/cutie_char_sprite.png"
    itemBoostImg.src = "media/speeditem.png"
    itemSlowImg.src = "media/slowitem.png"
    itemGhostImg.src = "media/ghostitem.png"
    playerGhostImg.src = "media/ghostPlayer.png"
    explosionImg.src = "media/explosion2.png"
    portalImg.src = "media/portal2.png"
    poisonImg.src = "media/poisonTest.png"
    info.append(nameLabel);
    info.append(posXLabel);
    info.append(posYLabel);
}

function drawImgXY (img, x,y, xh, yh, canvasX, canvasY, dw, dh) {

    drawX = x * playerImgWidth;
    drawY = y * playerImgHeight;

    ctx.drawImage(img, drawX, drawY, xh, yh,  canvasX, canvasY, dw, dh)

}

function drawPlayerChar (player, count) {

    let animationStateX = 1;

    if (count == frameLimit) {
        imgIndex++;
    }
    if (imgIndex == 4) {
        imgIndex = 0;
    }
    if (imgIndex == 0) {
        animationStateX = 1;
    }
    if(imgIndex == 1) {
        animationStateX = 0;
    }
    if(imgIndex == 2) {
        animationStateX = 1;
    }
    if(imgIndex == 3) {
        animationStateX = 2;
    }

    if (player.IsAlive && player.IsMoving){
        if (player.DirDown){
            drawImgXY(playerCharSprite, animationStateX, 0, playerImgWidth, playerImgHeight, player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if ( player.DirUp){
            drawImgXY(playerCharSprite, animationStateX, 1, playerImgWidth, playerImgHeight, player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if (player.DirLeft){
            drawImgXY(playerCharSprite, animationStateX, 2,  playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if (player.DirRight){
            drawImgXY(playerCharSprite, animationStateX, 3,  playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        }else {
            drawImgXY(playerCharSprite, animationStateX, 0, playerImgWidth, playerImgHeight, player.PositionX, player.PositionY, fieldSize, fieldSize)
        }
    } else if (player.IsAlive && !player.IsMoving) {
        if (player.DirDown){
            drawImgXY(playerCharSprite, 1, 0, playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if ( player.DirUp){
            drawImgXY(playerCharSprite, 1, 1, playerImgWidth, playerImgHeight, player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if (player.DirLeft){
            drawImgXY(playerCharSprite, 1, 2, playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        } else if (player.DirRight){
            drawImgXY(playerCharSprite, 1, 3, playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        }else {
            drawImgXY(playerCharSprite, 1, 1, playerImgWidth, playerImgHeight,player.PositionX, player.PositionY, fieldSize, fieldSize)
        }
    } else {
        ctx.drawImage(playerGhostImg, player.PositionX, player.PositionY, fieldSize, fieldSize);
    }
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
function drawPlayerPosClient() {
    if(currentUser == null){return}
    ctx.fillStyle = "rgba(0, 0, 0, 0.4)";
    const x = Math.floor((currentUser.PositionX + fieldSize/2)/fieldSize) * fieldSize
    const y = Math.floor((currentUser.PositionY + fieldSize/2)/fieldSize) * fieldSize
    ctx.fillRect(x,y , fieldSize , fieldSize)
}

function updateUI (sessionRunning) {
    if (sessionRunning && !exec) {
        exec = true;
        readyInfo.innerHTML = "Session is running.";
        startCountdown(countDown);
        console.log("session started");
    }
    if (!sessionRunning && exec) {
        clearInterval(intervall)
        exec = false;
        countDown.innerHTML = ""
        readyInfo.innerHTML = "<button id='readyButton' onclick='sendGetReady()'>Ready?</button>";
        readyButton = document.querySelector("#readyButton")
        console.log("session stopped")
    }
}

function startCountdown (container){
    let dateStart = new Date();
    let timeEnd = new Date().setMinutes(dateStart.getMinutes() + suddenDeathTimer)

    intervall = setInterval( function () {
        let timeNow = new Date().getTime();

        let timeDiff = timeEnd - timeNow;

        let minutes = Math.floor((timeDiff % (1000 * 60 * 60)) / (1000 * 60));
        let seconds = Math.floor((timeDiff % (1000 * 60)) / 1000);

        container.innerHTML = minutes + "m " + seconds + "s ";

        if (timeDiff < 0) {
            container.innerHTML =  "SUDDEN DEATH!";
        }

    })
}