
let isBombLegal = true;
let keyPresses = {};

document.addEventListener('keydown', keyDownListener, false);
document.addEventListener('keydown', spaceKeyDownListener, false);
document.addEventListener('keyup', keyUpListener, false);

function spaceKeyDownListener(event) {
    if(isBombLegal && event.key === " "){
        keyPresses[event.key] = true
        isBombLegal = false
        setTimeout(()=>{
            isBombLegal = true
        },bombTimeOutMS)
    }
}


function keyDownListener(event) {
    if(event.key !== " "){
        keyPresses[event.key] = true;
    }
}

function keyUpListener(event) {
    keyPresses[event.key] = false;
}
