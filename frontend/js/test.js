function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}


function postUploadImageRequest(){
    const inputFile = document.querySelector("#avatar").files[0];
    const request = createAjaxRequest();
    const formData = new FormData();
    formData.append("imageFile", inputFile);
    request.onreadystatechange = function(){
        if(this.readyState === 4){
            if(this.status === 200){

            }else{
                alert("failed: " +  this.status);
            }
        }
    }
    request.open("POST","/uploadImage",true);
    request.send(formData);
}

function openWelcome() {
    window.location.href="/index.html";
}