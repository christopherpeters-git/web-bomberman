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
            }else {
                readyButton.style.backgroundColor = 'green'
            }
        })
        .catch((reason => {
            console.log(reason)
        }))
}