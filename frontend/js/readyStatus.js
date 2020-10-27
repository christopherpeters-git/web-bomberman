function sendGetReady(){

    fetch("/setReady")
        .then(response => {
            if(response.status === 200){
                return response.text()
            }
            console.log(response.status)
        })
        .then(text => {
            if(text === "nrdy"){
                readyButton.style.backgroundColor = 'red'
                console.log("RDY clicked RED");
                console.log(readyButton)
            }else {
                readyButton.style.backgroundColor = 'green'
                console.log("RDY clicked GREEN");
                console.log(readyButton)
            }
        })
        .catch((reason => {
            console.log(reason)
        }))
}

