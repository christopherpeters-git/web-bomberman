let currentUser;

function updateUserInfo() {
    if(currentUser == null || nameLabel == null || posXLabel == null || posYLabel == null){return}
    nameLabel.innerHTML = "Name: " + currentUser.Name;
    posXLabel.innerHTML = "Position X: " + currentUser.PositionX;
    posYLabel.innerHTML = "Position Y: " + currentUser.PositionY;
}

function searchForUser(playerArr){
    if (userId === "" || playerArr == null){
        return
    }
    for(let i = 0; i < playerArr.length; i++) {
        if (playerArr[i].UserID == userId) {
            currentUser = playerArr[i]
            return
        }
    }
}