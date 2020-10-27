let loginActive = false;
let registerActive = false;

let loginButton = document.querySelector("#loginBtn");
let registerButton = document.querySelector("#regBtn");

function createAjaxRequest() {
    let request;
    if (window.XMLHttpRequest) {
        request = new XMLHttpRequest();
    } else {
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}
function openTab(evt, tabName) {
    let reg = document.querySelector("#regForm");
    let log = document.querySelector("#logForm");
    let todo = document.querySelector("#" + tabName);

    if (tabName == "regForm") {
        registerActive =true;
        loginActive = false;
    } else {
        registerActive = false;
        loginActive = true;
    }

    reg.style.display = "none";
    log.style.display = "none";

    todo.style.display = "initial";
}
function sendPostLoginRequest() {
    let username = document.getElementById("usernameLogin").value;
    let password = document.getElementById("passwordLogin").value;
    let container = document.querySelector(".loginFlexContainer");
    const request = createAjaxRequest();
    //todo check for illegal chars
    console.log(username + password);
    request.onreadystatechange = function () {
        if(4 === this.readyState) {
            container.innerHTML = "<h1> Das Hat geklappt!</h1> <h2>Du wirst nun weitergeleitet.</h2>"
            //todo fix browser back button
            //todo check if login was REALLY successful
            window.setTimeout(function redirect () {
                //-----must change link on pi-----
                window.location = "http://localhost:2100/game.html"
                container.innerHTML = "";
            },2000);

        }
    }
    request.open("POST","/login", true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("usernameInput="+username+"&"+"passwordInput="+password);
}

function sendPostRegisterRequest() {
    const username = document.getElementById("usernameRegister").value;
    const password = document.getElementById("passwordRegister").value;
    const request = createAjaxRequest();
    //todo check for illegal chars
    //todo check if register was REALLY successful and redirect if so
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            alert("Registrierung erfolgreich");
        }
    }
    request.open("POST","/register",true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("usernameInput="+username+"&"+"passwordInput="+password);
}

document.addEventListener('keypress', submitOnEnter, false);

function submitOnEnter (event) {
    if (event.keyCode === 13) {
        if (loginActive) {
            loginButton.click();
        }
        else if (registerActive) {
            registerButton.click();
        }
    }
}