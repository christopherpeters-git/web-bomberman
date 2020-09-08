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

    reg.style.display = "none";
    log.style.display = "none";

    todo.style.display = "initial";
}
function sendPostLoginRequest() {
    let username = document.getElementById("usernameLogin").value;
    let password = document.getElementById("passwordLogin").value;
    const request = createAjaxRequest();
    //todo check for illegal chars
    console.log(username + password);
    request.onreadystatechange = function () {
        if(4 === this.readyState) {
            alert("Login erfolgreich");
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
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            alert("Registrierung erfolgreich");
        }
    }
    request.open("POST","/register",true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("usernameInput="+username+"&"+"passwordInput="+password);
}