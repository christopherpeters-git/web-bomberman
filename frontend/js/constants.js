
const nameLabel = document.createElement("p");
const posXLabel = document.createElement("p");
const posYLabel = document.createElement("p");
const readyInfo = document.querySelector(".noShake")
const countDown = document.querySelector(".countdown");
const info = document.querySelector('#stats')
const ctx = document.querySelector('#matchfield').getContext("2d")
let readyButton = document.querySelector("#readyButton")

//Every *frameLimit* messages between Server and Client the playermodel changes, which animates the movement
const frameLimit = 8;

const fieldSize = 50;
const canvasSize = 1000;

const playerImgHeight = 32;
const playerImgWidth = 32;

//Time in which no Bomb can be placed after placing a Bomb
const bombTimeOutMS = 1000;

const suddenDeathTimer = 1;