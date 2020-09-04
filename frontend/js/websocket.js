let socket = new WebSocket("ws://localhost:80/ws-test/")
console.log("Attempting Websocket connection")

socket.onopen = () => {
    console.log("Connected")
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}

document.addEventListener( 'keydown', handleKeyPress, false );

function handleKeyPress(event){
    let keyCode;
    if (event.key !== undefined) {
        keyCode = event.key;
    } else if (event.keyIdentifier !== undefined) {
        keyCode = event.keyIdentifier;
    } else if (event.keyCode !== undefined) {
        keyCode = event.keyCode;
    }
    socket.send(keyCode);
}