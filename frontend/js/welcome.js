var slideIndex = 1;


// Next/previous controls
function plusSlides(n) {
    showSlides(slideIndex += n);
}

function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}

// Thumbnail image controls
function currentSlide(n) {
    showSlides(slideIndex = n);
}

function showSlides(n) {
    var i;
    var slides = document.getElementsByClassName("mySlides");
    var dots = document.getElementsByClassName("dot");
    if (n > slides.length) {slideIndex = 1}
    if (n < 1) {slideIndex = slides.length}
    for (i = 0; i < slides.length; i++) {
        slides[i].style.display = "none";
    }
    for (i = 0; i < dots.length; i++) {
        dots[i].className = dots[i].className.replace(" active", "");
    }
    slides[slideIndex-1].style.display = "block";
    dots[slideIndex-1].className += " active";
}

function initWelcome() {
    //toggleSlideshow();
}

function toggleSlideshow() {
    let slideshow = document.getElementById("slideshow");
    if(slideshow.hidden) {
        slideshow.hidden = false;
    }else{
        slideshow.hidden = true;
    }
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

function openTab(evt, tabName)  {
    // Declare all variables
    var i, tabcontent, tablinks;

    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the current tab, and add an "active" class to the button that opened the tab
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}