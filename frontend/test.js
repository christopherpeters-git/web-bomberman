function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}


function GetFetchNumberRequest(){
    const inputNumber = document.querySelector("#inputNumber").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function(){
        if(this.readyState === 4){
            if(this.status === 200){
                alert("Response: " + this.responseText);
            }else{
                alert("failed: " +  this.status);
            }
        }
    }
    request.open("GET","/fetchNumber?number="+inputNumber,true);
    request.send();
}