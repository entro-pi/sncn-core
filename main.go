package main

import (
  "time"
  "strconv"
  "strings"
  "fmt"
  "os"
  "io/ioutil"
  "github.com/SolarLune/dngn"
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  zmq "github.com/pebbe/zmq4"
//  "github.com/gorilla/websocket"
  "golang.org/x/net/websocket"
  "encoding/json"
)


type Descriptions struct {
	BATTLESPAM int
	ROOMDESC int
	PLAYERDESC int
	ROOMTITLE int
}
type Chat struct {
	User string
	Message string
	Time time.Time
}
type Space struct{
	Room dngn.Room
	Vnums string
	Zone string
	ZonePos []int
	ZoneMap [][]int
	Vnum int
	Desc string
	Mobiles []int
	Items []int
	CoreBoard string
	Exits Exit
	Altered bool
}
type Exit struct {
	North int
	South int
	East int
	West int
	NorthWest int
	NorthEast int
	SouthWest int
	SouthEast int
	Up int
	Down int
}

type Player struct {
	Name string
	Title string
	Inventory []int
	Equipment []int
	CoreBoard string
	PlainCoreBoard string
	CurrentRoom Space
  PlayerHash string

	MaxRezz int
	Rezz int
	Tech int

	Str int
	Int int
	Dex int
	Wis int
	Con int
	Cha int
}
type Payload struct {
  Client_id string
  Client_secret string
  Supports []string
  Channels []string
  Version string
  User_agent string
  Unicode string
  Status string

}
type HeartPayload struct {
  Players []string
}

type Heartbeat struct {
  Event string
  Payload HeartPayload
}

type Authenticator struct {
  Event string
  Payload Payload
}

type Mobile struct {
	Name string
	LongName string
	ItemSpawn []int
	Rep string
	MaxRezz int
	Rezz int
	Tech int
	Aggro int
	Align int
}



const (
	cmdPos = "\033[51;0H"
	mapPos = "\033[1;51H"
	descPos = "\033[0;50H"
	chatStart = "\033[38:2:200:50:50m{{=\033[38:2:150:50:150m"
	chatEnd = "\033[38:2:200:50:50m=}}"
	end = "\033[0m"
  grapevine = "wss://grapevine.haus/socket"
  callback = "https://snowcrashnetwork.grapevine.haus/auth"
)

// client dials the WebSocket echo server at the given url.
// It then sends it 5 different messages and echo's the server's
// response to each.
func client(clientid string, secret string, url string, player chan string, vineOn chan bool) error {
    _, cancel := context.WithTimeout(context.Background(), time.Minute)
    defer cancel()
    vineOn <- true
    //var playerList []string
    fmt.Println("IN CLIENT FUNC")
    var playing []string
    ws, err := websocket.Dial(url, "", callback)
    if err != nil {
        return err
    }
    defer ws.Close()
    var auth Authenticator
    auth.Event = "authenticate"
    auth.Payload.Client_id = strings.TrimSpace(clientid)
    auth.Payload.Client_secret = strings.TrimSpace(secret)
    auth.Payload.Supports = append(auth.Payload.Supports, "channels")
    auth.Payload.Channels = append(auth.Payload.Channels, "grapevine")
    auth.Payload.Version = "1.0.0"
    auth.Payload.User_agent = "snowcrashnetwork"
    authJSON, err := json.Marshal(auth)
    if err != nil {
      panic(err)
    }
    authJSONString := strings.ToLower(string(authJSON))
    fmt.Println(authJSONString)
    authorized := false
    //authSend := []byte(auth)
    for {
      if !authorized {
        //Authenticate
        _, err = ws.Write([]byte(authJSONString))
        if err != nil {
            return err
        }

        var msg = make([]byte, 1024)
        _, err = ws.Read(msg)
        if err != nil {
            return err
        }
        fmt.Println(string(msg))
        authorized = true
//        var authorizedMess Authenticator
//        err = json.Unmarshal(msg, &authorizedMess)
//        if err != nil {
//          fmt.Println(err)
//        }
//        if authorizedMess.Payload.Status == "success" {
//          authorized = true
//        }

      }else {


        var heart Heartbeat
        heart.Event = "heartbeat"
        select {
        case playersLog := <- player:
          enqueue := true
          logout := false
          playerName := ""
          if strings.Contains(playersLog, "LOGOUT||") {
            playerName = strings.Split(playersLog, "||")[1]
            logout = true
            enqueue = false
          }
          for i := 0;i < len(playing);i++ {
            if playing[i] == playersLog {
              enqueue = false
            }
          }
          if enqueue {
            if playersLog != "" {
              playing = append(playing, playersLog)
              heart.Payload.Players = append(heart.Payload.Players, playersLog)
            }
          }else if logout {
            var loggedOut []string
            for i := 0;i < len(playing);i++ {
              if strings.ToLower(playing[i]) == strings.ToLower(playerName) {
                if len(playing) > 1 {
                  continue
                }else if len(playing) == 1 {
                  loggedOut = append(loggedOut, "none")
                  playing = loggedOut
                }else {
                  loggedOut = append(loggedOut, playing[i])
                }

              }
              heart.Payload.Players = loggedOut
              logout = false
            }
          }else {
            heart.Payload.Players = playing
          }
        default:
          if len(heart.Payload.Players) <= 0  && len(playing) >= 1 {
            heart.Payload.Players = playing
          }else if len(heart.Payload.Players) <= 0 {
            heart.Payload.Players = append(heart.Payload.Players, "none")
          }
          if heart.Payload.Players[0] != "none" {
            heart.Payload.Players = playing
          }
        }
        heartJSON, err := json.Marshal(heart)
        if err != nil {
          panic(err)
        }
        heartJSONString := strings.ToLower(string(heartJSON))
        var heartbeat = make([]byte, 512)
        _, err = ws.Read(heartbeat)
        if err != nil {
            return err
        }
        fmt.Println(heartJSONString)
        fmt.Println(string(heartbeat))
        if strings.Contains(string(heartbeat), "heartbeat") {
          fmt.Println("\033[38:2:200:0:0mBeep.\033[0m")
          _, err = ws.Write([]byte(heartJSONString))
          if err != nil {
              return err
          }


        }



      }

    }
    vineOn <- false
    return nil
}
func grapeVine(playerList chan string, vineOn chan bool){
  clientFile, err := os.Open("client")
  if err != nil {
    panic(err)
  }
  clientid, err := ioutil.ReadAll(clientFile)
  secretFile, err := os.Open("secret")
  if err != nil {
    panic(err)
  }
  secret, err := ioutil.ReadAll(secretFile)
  if err != nil {
    panic(err)
  }
  fmt.Println(strings.TrimSpace(string(clientid)))
  fmt.Println(strings.TrimSpace(string(secret)))
  go client(string(clientid), string(secret), grapevine, playerList, vineOn)
}

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*32+100)
  }
  return newVal
}
func lookupPlayer(pass string) Player {
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  var player Player
  collection := client.Database("pfiles").Collection("Players")
  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":hash(pass)}})
  if err != nil {
    panic(err)
  }
  err = result.Decode(&player)
  if err != nil {
    fmt.Println("\033[38:2:150:0:150mPlayerfile requested was not found\033[0m")
    var noob Player
    noob.PlayerHash = "2"
    return noob
  }
  return player

}
func initPlayer(name string, pass string) Player {
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  play := InitPlayer(name, pass)
  playBSON, err := bson.Marshal(play)
  if err != nil {
    panic(err)
  }
  collection := client.Database("pfiles").Collection("Players")
  _, err = collection.InsertOne(context.Background(), playBSON)
  if err != nil {
    panic(err)
  }
  return play

}

func loopInput(servepubKey string, in chan string, players chan string, vineOn chan bool) {
  fmt.Println("Core login procedure started")

  response, err := zmq.NewSocket(zmq.REP)
  if err != nil {
    panic(err)
  }
  defer response.Close()
  hostname := "tcp://*:7777"
  err = response.Bind(hostname)
  if err != nil {
    panic(err)
  }
  var playerList []string
  for {
    select {
    case isVineOn := <- vineOn:
      if isVineOn == true {
        fmt.Println("Grapevine \033[38:2:0:200:0mActive\033[0m")
      }
      if isVineOn == false {
        fmt.Println("Grapevine \033[38:2:200:0:0mInactive\033[0m")
        grapeVine(players, vineOn)
        time.Sleep(15*time.Second)
      }
    default:
      fmt.Println("Grapevine Capable")
    }
    fmt.Println("IN LOOP")
    request, err := response.Recv(0)
    if err != nil {
      panic(err)
    }
    var play Player
    in <- request
    if strings.Contains(request, ":"){
      fmt.Println(strings.Split(request, ":")[0]+"AUTHINFO")
    }else if strings.Contains(request, "+=+") {
      fmt.Println(strings.Split(request, "+=+")[0]+"CHATMSGS")
    }else {
      fmt.Println(request)
    }
    if strings.Contains(request, ":-:") {
        userPass := strings.Split(request, ":-:")
        name, pass := userPass[0], userPass[1]
        play = initPlayer(name, pass)
        _, err = response.Send(play.PlayerHash, 0)
        if err != nil {
          panic(err)
        }
    }else if strings.Contains(request, ":=:") {
        userPass := strings.Split(request, ":=:")
        pass := userPass[1]
        play = lookupPlayer(pass)
        playBytes, err := bson.Marshal(play)
        if err != nil {
          panic(err)
        }
        fmt.Println(play.Name+"LOGGED IN")
        playerList = append(playerList, play.Name)
        _, err = response.SendBytes(playBytes, 0)
        if err != nil {
          panic(err)
        }
        players <- play.Name
    }else if strings.Contains(request, "+=+") {
      message := strings.Split(request, "+=+")[1]
      playerName := strings.Split(request, "+=+")[0]
      fmt.Println("Creating chat")
      createChat(message, playerName)
      toSend := showChat()
      fmt.Println("Sending chat")
      _, err = response.Send(toSend, 0)
      if err != nil {
        panic(err)
      }

    }
    request, err = response.Recv(0)
    if err != nil {
      panic(err)
    }
    if strings.Contains(request, "+==LOGOUT") {
        playerToLogOut := strings.Split(request, "+==LOGOUT")[0]
        for i := 0;i < len(playerList);i++ {
          if strings.ToLower(playerList[i]) == strings.ToLower(playerToLogOut) {
            fmt.Println(playerList[i]+ "LOGGED OUT")
            response.Send("LOGGED "+playerList[i]+" OUT", 0)
            playerList[i] = ""
            players <- "LOGOUT||"+strings.ToLower(playerToLogOut)

          }
        }
    }else if request == "+===shutdown===+" {

      fmt.Println("\033[38:2:255:0:0mGOT "+request+" SIGNAL\033[0m")
      os.Exit(1)
    }else if strings.Contains(request, ":go to=") {
      if len(strings.Split(request, ":")) == 2 {
    //    playerHash := strings.Split(request, ":")[0]
    //    gotovnum := strings.Split(request, "=")[1]
  //      play = lookupPlayer(pass)
//        goTo(dest int, play Player, populated []Space)
      }
    }else {

  //    in <- request
      _, err := response.Send("INVALID REQUEST", 0)
      fmt.Println("\033[38:2:150:0:150m"+request+"\033[0m")
      if err != nil {
        panic(err)
      }
    }
  }

}

func main() {
    vineOn := make(chan bool)
    in := make(chan string)
    var play Player
    var populated []Space
    playerList := make(chan string)
    grapeVine(playerList, vineOn)
    clientkey, _, err := zmq.NewCurveKeypair()
    if err != nil {
      panic(err)
    }
    go loopInput(clientkey, in, playerList, vineOn)
    for {
      value := <-in
      if strings.HasPrefix(value, "init world:") {
        playerName := strings.Split(value, "ld:")[1]
        pass := strings.Split(value, "--")[1]
        descString := "The absence of light is blinding.\nThree large telephone poles illuminate a small square."
  			for len(strings.Split(descString, "\n")) < 8 {
  				descString += "\n"
  			}
  			InitZoneSpaces("0-5", "The Void", descString)
  			descString = "I wonder what day is recycling day.\nEven the gods create trash."
  			for len(strings.Split(descString, "\n")) < 8 {
  				descString += "\n"
  			}
  			InitZoneSpaces("5-15", "Midgaard", descString)
  			populated = PopulateAreas()
        play = InitPlayer(playerName, pass)
        play = InitPlayer("dorp", "norp")

  			addPfile(play)
  			createMobiles("Noodles")
        respond := fmt.Sprint("\033[38:2:0:250:0mInitialized "+strconv.Itoa(len(populated))+" rooms\033[0m")
        fmt.Printf(respond)

        fmt.Println("\033[38:2:0:250:0mAll tests passed and world has been initialzed\n\033[0mYou may now start with --login.")

      }
    }
}
