var stage = document.getElementById("matchfield"), // Get the canvas element by Id
    ctx = stage.getContext("2d"),           // Canvas 2d rendering context
    x = 10,                                         //intial horizontal position of Player Figure
    y = 10,                                         //intial vertical position of of Player Figure
    wid = 50,                                       //width of of Player Figure
    hei = 50;                                       //height of of Player Figure

var img = document.getElementById("testImg");

//Draw Fig function
function drawFigure(x, y, wid, hei) {
    ctx.drawImage(img, x, y, wid, hei);
}

drawFigure(x, y, wid, hei); //Drawing of Player Figure on initial load

window.onkeydown = function (event) {
    var keyPr = event.keyCode; //Key code of key pressed

    if (keyPr === 39 && x <= 460) {
        x = x + 20; //right arrow add 20
        console.log("RIGHT")
    } else if (keyPr === 37 && x > 10) {
        x = x - 20; //left arrow subtract 20
        console.log("LEFT")
    } else if (keyPr === 38 && y > 10) {
        y = y - 20; //top arrow subtract 20
        console.log("TOP")
    } else if (keyPr === 40 && y <= 460) {
        y = y + 20; //down arrow add 20
        console.log("DOWN")
    }


    ctx.clearRect(0, 0, 500, 500);

    //Drawing playerFig at new position
    drawFigure(x, y, wid, hei);
};