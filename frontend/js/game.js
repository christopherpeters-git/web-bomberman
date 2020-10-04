

function gameLoop (map, players){
    //Needed for movement Animation in display.js
    frameCounter++;

    //Updates user Information for logged in User
    searchForUser(players)
    updateUserInfo()

    ctx.clearRect(0, 0, canvasSize, canvasSize);
    background(grassImg, map);
    drawPlayerPosClient();

    //Every Player and their name is drawn on the canvas
    for(let i = 0; i < players.length; i++){
        ctx.font = "normal 8px 'Press Start 2P'";
        ctx.textAlign= "center";
        ctx.fillStyle = "rgba(0,0,0,0.75)"
        ctx.fillText(players[i].Name,players[i].PositionX + 25,players[i].PositionY - 5, 100);
        if (players[i] != null){
            if (!players[i].GhostActive){
                drawPlayerChar(players[i], frameCounter)
            }
        }
    }
    //FildObjects are Drawn onto the Canvas
    drawImageFromEnum(itemBoostImg, map, 6);
    drawImageFromEnum(itemSlowImg, map, 7);
    drawImageFromEnum(itemGhostImg, map, 8);
    drawImageFromEnum(portalImg, map, 12)
    drawImageFromEnum(wallImg, map, 3);
    drawImageFromEnum(wallImg2, map, 2);
    drawImageFromEnum(bomb3Img, map, 1);
    drawImageFromEnum(bombImg, map, 10);
    drawImageFromEnum(bomb2Img, map, 11);
    drawImageFromEnum(explosionImg, map, 9)
    drawImageFromEnum(poisonImg, map, 13)

    //When GhostActive is true, players are drawn after other objects, so that they "float" above them.
    for(let i = 0; i < players.length; i++){
        if(players[i] != null){
            if (players[i].GhostActive){
                ctx.globalAlpha = 0.5
                drawPlayerChar(players[i], frameCounter)
                ctx.globalAlpha = 1
            }
        }
    }
    //Needed for movement Animation in display.js
    if (frameCounter == frameLimit){
        frameCounter = 0;
    }
}

