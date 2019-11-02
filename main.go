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
  grapevine = "tcp://grapevine.haus/socket"
  callback = "tcp://snowcrashnetwork.vineyard.haus/auth"
)

func grapeVine() {
  sendSocket, err := zmq.NewSocket(zmq.REQ)
  if err != nil {
    panic(err)
  }

  response, err := zmq.NewSocket(zmq.REP)
  if err != nil {
    panic(err)
  }
  err = response.Connect("tcp://*:7787")
  if err != nil {
    panic(err)
  }

  clientFile, err := os.Open("client")
  if err != nil {
    panic(err)
  }
  clientid, _ := ioutil.ReadAll(clientFile)
  secretFile, err := os.Open("secret")
  if err != nil {
    panic(err)
  }
  secret, _ := ioutil.ReadAll(secretFile)


  sendSocket.Connect(grapevine)
  auth := `{  "event": "authenticate",  "payload": {    "client_id": "`+string(clientid)+`",    "client_secret": "`+string(secret)+`",    "supports": ["channels"],    "channels": ["grapevine"],    "version": "1.0.0",    "user_agent": "snowcrash.network v 0.01"  }}`
  _, err = sendSocket.Send(auth, 0)

  for {
    input, err := sendSocket.Recv(0)
    if err != nil {
      panic(err)
    }
    fmt.Println(input)
    result, err := response.Recv(0)
    if err != nil {
      panic(err)
    }
    fmt.Println(result)
    os.Exit(1)
  }
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

func loopInput(servepubKey string, in chan string) {
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
  for {
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

    if strings.Split(string(request), ":")[0] == "REQUESTPUBKEY" {

        _, err = response.Send(servepubKey, 0)

        if err != nil {
          panic(err)
        }
    }else if strings.Contains(request, ":-:") {
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

        _, err = response.SendBytes(playBytes, 0)
        if err != nil {
          panic(err)
        }
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
      if err != nil {
        panic(err)
      }
    }
  }

}

func main() {
    in := make(chan string)
    var play Player
    var populated []Space
    grapeVine()
    clientkey, _, err := zmq.NewCurveKeypair()
    if err != nil {
      panic(err)
    }
    go loopInput(clientkey, in)
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

      }else if value == "shutdown" {

        fmt.Println("\033[38:2:255:0:0mGOT "+value+" SIGNAL\033[0m")
        os.Exit(1)
      }
    }
}
