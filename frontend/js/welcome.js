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


function sendPostLoginRequest() {}

function sendPostRegisterRequest() {}
