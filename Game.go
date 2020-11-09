package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

//Constant
const (
	FIELD_SIZE                  = 50
	STEP_SIZE                   = 3
	CANVAS_SIZE                 = 1000
	STANDARD_BOMB_RADIUS        = 2
	STANDARD_BOMB_TIME          = 3
	STANDARD_STEP_MULTIPLICATOR = 1
	DEATH_STEP_MULTIPLICATOR    = 0.5
	SUDDEN_DEATH_START_TIME     = 60
	MAP_SIZE                    = CANVAS_SIZE / FIELD_SIZE

	/*
		10 is equal to full map, 10 is MAX!!!
	*/
	SUDDEN_DEATH_MAX_AREA = 7

	/*
		In seconds, higher number means more time between the increase of the area
	*/
	SUDDEN_INCREASE_TIME = 5
)

var GameMap Map
var Connections map[uint64]*Session
var ticker *time.Ticker
var spawnPositions [][]int

//var incomingTicker = time.NewTicker(1 * time.Millisecond)
var sessionRunning bool
var suddenDeathRunning bool

/*Things send to the clients*/
var bombermanArray []Bomberman
var abstractGameMap chan [][][]FieldObject
var clientPackageAsJson []byte

//Called before any connection is possible
func initGame() {
	//Global variables
	GameMap = NewMap(MAP_SIZE)
	Connections = make(map[uint64]*Session, 0)
	ticker = time.NewTicker(16 * time.Millisecond)
	spawnPositions = [][]int{{0, 0}, {0, 10}, {0, 19}, {10, 0}, {10, 19}, {19, 0}, {19, 10}, {19, 19}}
	sessionRunning = false
	suddenDeathRunning = false
	bombermanArray = make([]Bomberman, 0)
	//abstractGameMap = make([][][]FieldObject,0)
	abstractGameMap = make(chan [][][]FieldObject)
	clientPackageAsJson = make([]byte, 0)

	//Routines
	go UpdateClients()
	go func() {
		for {
			abstractGameMap <- BuildAbstractGameMap()
		}
	}()
}

/*
Wrapper Function to Build the Map new.
*/
//func //MapChanged() {
//	go BuildAbstractGameMap()
//}

/*
Represents the Position of a Player.
*/
type Position struct {
	x int
	y int
}

/*
Initialises a new Position with x and y.
*/
func newPosition(x int, y int) Position {
	return Position{
		x: x,
		y: y,
	}
}

//Updates a Position-object
func (p *Position) updatePosition(xOffset int, yOffset int) {
	p.x += xOffset
	p.y += yOffset
}

/*
Converts a Pixelposition to an Arrayposition.
*/
func pixToArr(pixel int) int {
	return (pixel + FIELD_SIZE/2) / FIELD_SIZE
}

/*
Prints a list to the log.
*/
func printList(list *list.List) {
	element := list.Front()
	if element == nil {
		log.Println("List is null!")
		return
	}
	log.Println("List started: ")
	log.Println(element.Value.(*Bomberman))
	for element.Next() != nil {
		log.Println(element.Value.(*Bomberman))
		element = element.Next()
	}
	log.Println("List ended...")
}

/*
Represents the Keys pressed by the Player.
*/
type KeyInput struct {
	Wpressed     bool `json:"w"`
	Spressed     bool `json:"s"`
	Apressed     bool `json:"a"`
	Dpressed     bool `json:"d"`
	SpacePressed bool `json:" "`
}

/*
Information which the Client needs.
This things will be send to Client.
*/
type ClientPackage struct {
	Players        []Bomberman
	GameMap        [][][]FieldObject
	SessionRunning bool
}

/*
Wrapper for the user
*/
type Session struct {
	User              *User           //Connected user
	Bomber            *Bomberman      //Bomber of the connected user
	Connection        *websocket.Conn //Websocket connection
	ConnectionStarted time.Time       //point when player joined
}

func NewSession(user *User, character *Bomberman, connection *websocket.Conn, connectionStarted time.Time) *Session {
	return &Session{User: user, Bomber: character, Connection: connection, ConnectionStarted: connectionStarted}
}

/*
Returns the string representation of the connection
*/
func (r *Session) String() string {
	return "Session: { " + r.User.String() + "|" + r.Bomber.String() + "|" + r.Connection.RemoteAddr().String() + "|" + r.ConnectionStarted.String() + "}"
}

/*
Prints every active connection
*/
func AllConnectionsAsString() string {
	result := "Active Connections:"
	for _, v := range Connections {
		result += v.String() + "\n"
	}
	return result
}

/*
Starts the interaction loop.
*/
func StartPlayerLoop(session *Session) {
	//Add the infos to the connection map
	//MapChanged()
	Connections[session.User.UserID] = session
	playerWebsocketLoop(session)
	//Remove player from list at his last array position
	x := (session.Bomber.PositionX + FIELD_SIZE/2) / FIELD_SIZE
	y := (session.Bomber.PositionY + FIELD_SIZE/2) / FIELD_SIZE
	removePlayerFromList(GameMap.Fields[x][y].Player, session.Bomber)
	//Remove from the connection map
	delete(Connections, session.User.UserID)
}

/*
Interaction loop.
The user input is received and the Player-Position is updated accordingly / Bombs get placed.
*/
func playerWebsocketLoop(session *Session) {
	for {
		session.Bomber.IsMoving = false
		_, p, err := session.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var keys KeyInput
		if err := json.Unmarshal(p, &keys); err != nil {
			log.Println(err)
			continue
		}

		/*
			Checks which Key got pressed and performs an Action accordingly. If a movement key was pressed, the Collision and "legalness" of the Movement
			will be checked before updating the Player-Position.
		*/

		realStepSize := STEP_SIZE * session.Bomber.stepMult
		if keys.Wpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = true, false, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, -int(realStepSize)) {
				if session.Bomber.moveIfLegal(session.Bomber.PositionX, session.Bomber.PositionY-int(realStepSize)) {

					session.Bomber.topRightPos.updatePosition(0, -int(realStepSize))
					session.Bomber.topLeftPos.updatePosition(0, -int(realStepSize))
					session.Bomber.bottomRightPos.updatePosition(0, -int(realStepSize))
					session.Bomber.bottomLeftPos.updatePosition(0, -int(realStepSize))
					session.Bomber.PositionY -= int(realStepSize)
				}
			}
		} else
		//S
		if keys.Spressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, true, false, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(0, int(realStepSize)) {
				if session.Bomber.moveIfLegal(session.Bomber.PositionX, session.Bomber.PositionY+int(realStepSize)) {

					session.Bomber.topRightPos.updatePosition(0, int(realStepSize))
					session.Bomber.topLeftPos.updatePosition(0, int(realStepSize))
					session.Bomber.bottomRightPos.updatePosition(0, int(realStepSize))
					session.Bomber.bottomLeftPos.updatePosition(0, int(realStepSize))
					session.Bomber.PositionY += int(realStepSize)
				}
			}
		} else
		//A
		if keys.Apressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, true, false
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(-int(realStepSize), 0) {
				if session.Bomber.moveIfLegal(session.Bomber.PositionX-int(realStepSize), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.topLeftPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.bottomRightPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.bottomLeftPos.updatePosition(-int(realStepSize), 0)
					session.Bomber.PositionX -= int(realStepSize)
				}
			}
		} else
		//D
		if keys.Dpressed {
			session.Bomber.DirUp, session.Bomber.DirDown, session.Bomber.DirLeft, session.Bomber.DirRight = false, false, false, true
			session.Bomber.IsMoving = true
			if session.Bomber.collisionWithSurroundings(int(realStepSize), 0) {
				if session.Bomber.moveIfLegal(session.Bomber.PositionX+int(realStepSize), session.Bomber.PositionY) {

					session.Bomber.topRightPos.updatePosition(int(realStepSize), 0)
					session.Bomber.topLeftPos.updatePosition(int(realStepSize), 0)
					session.Bomber.bottomRightPos.updatePosition(int(realStepSize), 0)
					session.Bomber.bottomLeftPos.updatePosition(int(realStepSize), 0)
					session.Bomber.PositionX += int(realStepSize)
				}
			}
		}
		//Spacebar
		if keys.SpacePressed {
			if session.Bomber.IsAlive {
				go session.Bomber.placeBomb()
			}
		}
	}

}

/*
Updates the Client in an Interval.
*/
func UpdateClients() {
	for _ = range ticker.C {
		err := sendDataToClients()
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println("Updating Clients stopped.")
}

/*
Sends the Data needed by the Client to the Client.
*/
func sendDataToClients() error {
	//Create array from all connected Bombermen
	connectionLength := len(Connections)
	bombermanArray = make([]Bomberman, connectionLength)
	count := 0

	wg := &sync.WaitGroup{}
	wg.Add(connectionLength)
	for _, v := range Connections {
		session := v
		go func(count int) {
			bombermanArray[count] = *session.Bomber
			wg.Done()
		}(count)
		count++
	}
	wg.Wait()

	var err error
	clientPackageAsJson, err = json.Marshal(ClientPackage{
		Players:        bombermanArray,
		GameMap:        <-abstractGameMap,
		SessionRunning: sessionRunning,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range Connections {
		if err := v.Connection.WriteMessage(websocket.TextMessage, clientPackageAsJson); err != nil {
			return err
		}
	}
	return nil
}

func isLesserThan(a interface{}, b interface{}) bool {
	return a.(*Session).User.UserID < b.(*Session).User.UserID
}

/*
Starts the a Game-Session, if more then one Player is connected and all are ready.
*/
func StartGameIfPlayersReady() {
	if len(Connections) < 2 {
		return
	}
	for _, v := range Connections {
		if !v.Bomber.PlayerReady {
			return
		}
	}
	resetGame("images/testMap.png")
	sessionRunning = true
	for _, v := range Connections {
		v.Bomber.PlayerReady = false
	}
	time.AfterFunc(time.Second*SUDDEN_DEATH_START_TIME, startSuddenDeath)
}

/*
Starts the Suddendeath and Poison spreading.
*/
func startSuddenDeath() {
	suddenDeathRunning = true
	p := newPoison()
	//go checkForPoison()
	for t := 0; t < SUDDEN_DEATH_MAX_AREA; t++ {
		if !suddenDeathRunning {
			break
		}
		for i := 0; i < len(GameMap.Fields); i++ {
			for j := 0; j < len(GameMap.Fields[i]); j++ {
				if (i == t) || (j == t) || (i == 19-t) || (j == 19-t) {
					if GameMap.Fields[i][j].Contains[0] != nil {
						if GameMap.Fields[i][j].Contains[0].getType() == 13 {
							continue
						}
					}
					if GameMap.Fields[i][j].Contains[1] != nil {
						if GameMap.Fields[i][j].Contains[1].getType() == 13 {
							continue
						}
					}
					GameMap.Fields[i][j].addPoison(&p)
					killAllPlayersOnField(GameMap.Fields[i][j].Player)
				}
			}
		}
		findWinner()
		time.Sleep(time.Second * SUDDEN_INCREASE_TIME)
	}
}

/*
Inefficient! todo: Change!
While Sudden Death is running, constantly loops to all Fields and, if a Poison-Field is found, kills all player on the Field.
*/
//func checkForPoison() {
//	for suddenDeathRunning {
//		for i := 0; i < len(GameMap.Fields); i++ {
//			for j := 0; j < len(GameMap.Fields[i]); j++ {
//				if GameMap.Fields[i][j].Contains[0] != nil || GameMap.Fields[i][j].Contains[1] != nil{
//					if GameMap.Fields[i][j].Contains[0].getType() == 13 || GameMap.Fields[i][j].Contains[1].getType() == 13 {
//						if GameMap.Fields[i][j].Player != nil {
//							//TO DO: Dont insta kill
//							killAllPlayersOnField(GameMap.Fields[i][j].Player)
//							findWinner()
//						}
//					}
//				}
//			}
//		}
//	}
//}

/*
Resets the Game.
*/
func resetGame(s string) {
	suddenDeathRunning = false
	playerDied = false
	GameMap.clear()
	if err := CreateMapFromImage(GameMap, s); err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, v := range Connections {
		if count > 7 {
			count = 0
		}
		v.Bomber.Reset(spawnPositions[count][0], spawnPositions[count][1])
		count++
	}
}

/*
Kills all Players on a Field.
*/
func killAllPlayersOnField(list *list.List) {
	element := list.Front()
	if element != nil {
		element.Value.(*Bomberman).Kill()
		playerDied = true
		for element.Next() != nil {
			element = element.Next()
			element.Value.(*Bomberman).Kill()
		}
	}
}

/*
Checks if only one Player is alive and acts accordingly.
*/
func findWinner() {
	counter := 0
	var lastBomberAlive *Bomberman
	for _, v := range Connections {
		if v.Bomber.IsAlive {
			lastBomberAlive = v.Bomber
			counter++
		}
	}
	if counter > 1 {
		return
	} else if counter == 0 {
		log.Println("Draw")
	} else if counter == 1 {
		log.Println(lastBomberAlive.Name + "has Won")
		user, err := getUserByID(db, lastBomberAlive.UserID)
		if err != nil {
			log.Println(err)
		}
		user.GamesWon = user.GamesWon + 1

		err = updatePlayerStats(db, *user)
		if err != nil {
			log.Println(err)
		}
	}
	//todo send message
	resetGame("images/testMap.png")
	sessionRunning = false
}
