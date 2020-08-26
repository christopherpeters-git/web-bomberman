let socket = new WebSocket("ws://localhost:80/ws-test/")
console.log("Attempting Websocket connection")

socket.onopen = () => {
    console.log("Connected")
    socket.send("Hi from the client!")
}

socket.onclose = (event) => {
    console.log("Disconnected: " + event)
}

socket.onerror = (error) => {
    console.log("Connection failed: " + error)
}