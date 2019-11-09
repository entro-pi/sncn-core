package main

import (
  "time"
  "math/rand"
  "strconv"
  "strings"
  "fmt"
  "os"
  "io/ioutil"
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  zmq "github.com/pebbe/zmq4"
//  "github.com/gorilla/websocket"
  "golang.org/x/net/websocket"
  "encoding/json"
)



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
func client(broadcast []Broadcast, clientid string, secret string, url string, player chan string, vineOn chan bool, broadcastLine chan Broadcast) error {
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
          }else if strings.Contains(playersLog, "+|+") {
            fmt.Println("\033[38:2:255:0:0mTriggered subscribe\033[0m")
              channelSub := strings.Split(playersLog, "+|+")[1]

              //player should be assigned to this
              _ = strings.Split(playersLog, "+|+")[0]
              var ChannelSub GrapeMess
              ChannelSub.Event = "channels/subscribe"
              ChannelSub.Payload.Channel = channelSub
              ChannelSub.Ref = UIDMaker()

              ChannelSubJSON, err := json.Marshal(ChannelSub)
              if err != nil {
                panic(err)
              }
              ChannelSubJSONToSend := strings.ToLower(string(ChannelSubJSON))
              fmt.Println("\033[38:2:200:0:0mNEW SUBSCRIPTION.\033[0m")
              _, err = ws.Write([]byte(ChannelSubJSONToSend))
              if err != nil {
                  return err
              }
          continue
          }else if strings.Contains(playersLog, "||UWU||") {
            fmt.Println("\033[38:2:0:200:0mBroadcast Send\033[0m")
              channel := strings.Split(playersLog, "||UWU||")[1]
              playName := strings.Split(playersLog, "||UWU||")[0]
              messageParts := strings.Split(playersLog, "||}}{{||")[1]
              message := strings.Split(messageParts, "+++")[0]
              messageLong := strings.Split(messageParts, "+++")[1]
              //player should be assigned to this
              _ = strings.Split(playersLog, "||UWU||")[0]
              var Send SendBroadcast
              var Save Broadcast
              Save.Event = "channels/send"
              Send.Event = "channels/send"
              Save.Payload.Channel = channel
              Send.Payload.Channel = channel
              Send.Ref = UIDMaker()
              Save.Ref = Send.Ref
              Send.Payload.Name = playName
              Save.Payload.Name = playName
              Send.Payload.Message = message
              Save.Payload.Message = message
              Send.Payload.BigMessage = messageLong
              Save.Payload.BigMessage = messageLong
              SendJSON, err := json.Marshal(Send)
              initGrape(Save)
              if err != nil {
                panic(err)
              }
              SendJSONTo := strings.ToLower(string(SendJSON))
              fmt.Println("\033[38:2:0:200:0mNEW BROADSIDED BROADCAST.\033[0m")
              _, err = ws.Write([]byte(SendJSONTo))
              if err != nil {
                  return err
              }
          continue
          }else if strings.Contains(playersLog, "=+=") {
            var broadcastNow Broadcast
              broadcastLine <- broadcastNow
              continue
          }else if strings.Contains(playersLog, "-|-") {
            fmt.Println("\033[38:2:255:0:0mTriggered unsubscribe\033[0m")
              channelSub := strings.Split(playersLog, "-|-")[1]

              //player should be assigned to this
              _ = strings.Split(playersLog, "-|-")[0]
              var ChannelSub GrapeMess
              ChannelSub.Event = "channels/unsubscribe"
              ChannelSub.Payload.Channel = channelSub
              ChannelSub.Ref = UIDMaker()

              ChannelSubJSON, err := json.Marshal(ChannelSub)
              if err != nil {
                panic(err)
              }
              ChannelSubJSONToSend := strings.ToLower(string(ChannelSubJSON))
              fmt.Println("\033[38:2:200:0:0mNEW UNSUBSCRIPTION.\033[0m")
              _, err = ws.Write([]byte(ChannelSubJSONToSend))
              if err != nil {
                  return err
              }
          continue
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
                  if len(loggedOut) < 1 {
                    emptyWho := make([]string, 0)
                    loggedOut = emptyWho
                  }
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
            emptyWho := make([]string, 0)
            heart.Payload.Players = emptyWho
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
        if strings.Contains(string(heartbeat), "channels/broadcast") {
          fmt.Println("\033[38:2:0:150:150m[[Message]]\033[0m")
          var broadcastNow Broadcast
          newBroad := strings.Trim(string(heartbeat), "\x00")

          err = json.Unmarshal([]byte(newBroad), &broadcastNow)
          if err != nil {
            fmt.Println(err)
          }
          broadcast = append(broadcast, broadcastNow)
          broadcastLine <- broadcastNow
//          player <- fmt.Sprint("=+=")

          fmt.Println(broadcastNow.Payload.Name+"@"+broadcastNow.Payload.Game+"\033[38:2:0:150:150m[["+broadcastNow.Payload.Message+"]]\033[0m")
        }



      }

    }
    vineOn <- false
    return nil
}

func grapeVine(broadcast []Broadcast, playerList chan string, vineOn chan bool, broadcastLine chan Broadcast){
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
  go client(broadcast, string(clientid), string(secret), grapevine, playerList, vineOn, broadcastLine)
}

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*32+100)
  }
  return newVal
}
func updatePlayerSlain(hash string) {
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

  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":hash}})
  if err != nil {
    panic(err)
  }
  err = result.Decode(&player)
  if err != nil {
    fmt.Println("\033[38:2:150:0:150mPlayerfile requested was not found\033[0m")
  }

  player.Slain++

  if err != nil {
    panic(err)
  }
  _, err = collection.UpdateOne(context.Background(), options.Update().SetUpsert(true), bson.M{"$set":bson.M{"name":player.Name,"title":player.Title,"inventory":player.Inventory, "equipment":player.Equipment,
						"coreboard": player.CoreBoard,"currentroom":player.CurrentRoom,"slain":player.Slain, "hoarded":player.Hoarded, "str": player.Str, "int": player.Int, "dex": player.Dex, "wis": player.Wis, "con":player.Con, "cha":player.Cha, "classes": player.Classes }})

  if err != nil {
    panic(err)
  }







}
func lookupPlayer(name string, pass string) Player {
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

  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":hash(name+pass)}})
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
func getGrapes() []Broadcast {
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }

  collection := client.Database("broadcasts").Collection("general")
  curs, err := collection.Find(context.Background(), bson.M{})
  if err != nil {
    panic(err)
  }
  var broadContainer []Broadcast
  for curs.Next(context.Background()) {
    var broad Broadcast
    err = curs.Decode(&broad)
    if err != nil {
      panic(err)
    }
    broadContainer = append(broadContainer, broad)
  }

  return broadContainer

}
func initGrape(bcast Broadcast) Broadcast {
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  update := bson.M{"event":bcast.Event,"ref":bcast.Ref,"payload":bson.M{"channel":bcast.Payload.Channel, "message":bcast.Payload.Message,"game":bcast.Payload.Game,"name":bcast.Payload.Name}}
  collection := client.Database("broadcasts").Collection("snowcrash")
  _, err = collection.InsertOne(context.Background(), update)
  if err != nil {
    panic(err)
  }
  fmt.Println("Upserted the broadcast")
  return bcast

}

func loopInput(populated []Space, broadcast []Broadcast, in chan string, players chan string, vineOn chan bool, broadcastLine chan Broadcast) {
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
//  var broadcastContainer []Broadcast
  broadcastContainer := getBroadcasts()
  var playerList []Player
  for {
    select {
    case isVineOn := <- vineOn:
      if isVineOn == true {
        fmt.Println("Grapevine \033[38:2:0:200:0mActive\033[0m")
      }
      if isVineOn == false {
        fmt.Println("Grapevine \033[38:2:200:0:0mInactive\033[0m")
        grapeVine(broadcast, players, vineOn, broadcastLine)
        time.Sleep(15*time.Second)
      }
    case broadShip := <- broadcastLine:
      fmt.Println("Triggered broadSideLine")
      broadcast = getGrapes()
      broadcast = append(broadcast, initGrape(broadShip))
      if err != nil {
        panic(err)
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

    if strings.Contains(request, "\"event\": \"restart\",") {
      time.Sleep(15*time.Second)
    }
    if strings.Contains(request, "+|+") {

        fmt.Println("Starting the grapeclient")
        players <- request

    }
    if strings.Contains(request, "-|-"){
      fmt.Println("Starting the grapeclient")
      players <- request
    }
    if strings.Contains(request, "||UWU||") {
      fmt.Println("GRAPE BROADSIDE")
      players <- request
    }
    if strings.Contains(request, ":-:") {
        userPass := strings.Split(request, ":-:")
        name, pass := userPass[0], userPass[1]
        play = InitPlayer(name, pass)
        play.CurrentRoom = populated[0]
        addPfile(play, pass)
        savePfile(play)
        playBytes, err := bson.Marshal(play)
        _, err = response.SendBytes(playBytes, 0)
        if err != nil {
          panic(err)
        }
    }else if strings.Contains(request, ":=:") {
        userPass := strings.Split(request, ":=:")
        pass := userPass[1]
        name := userPass[0]
        play = lookupPlayer(name, pass)
        fmt.Println(play.Name+"LOGGED IN")
        play.Session = UIDMaker()
        playerList = append(playerList, play)

        playBytes, err := bson.Marshal(play)
        if err != nil {
          panic(err)
        }
        _, err = response.SendBytes(playBytes, 0)
        if err != nil {
          panic(err)
        }
        players <- play.Name
    }else if strings.Contains(request, "::CHECK::") {
        playerHash := strings.Split(request, "::CHECK::")[0]
        playerSession := strings.Split(request, "::CHECK::")[1]
        done := false
        for i := 0;i < len(playerList);i++ {
          if playerList[i].PlayerHash == playerHash {
            if playerList[i].Session == playerSession {
              fmt.Println(playerList[i].Session)
              _, err  = response.Send("OK", 0)
              if err != nil {
                panic(err)
              }
              done = true
            }
          }
        }
        if !done {
          _, err = response.Send("+__+SHUTDOWN+__+", 0)
          if err != nil {
            fmt.Println(err)
          }
        }
      }else if request == "::INVALIDATE::" {
        for i := 0;i < len(playerList);i++ {
          playerList[i].Session = "INVALID SELECTION"
        }
        _, err = response.Send("DONE", 0)
        if err != nil {
          panic(err)
        }
    }else if strings.Contains(request, "+++") {
      toSend := showChat()
      fmt.Println("Sending chat list")
      _, err = response.Send(toSend, 0)
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

    }else if strings.Contains(request, "++SAVE++"){
      if len(strings.Split(request, "++SAVE++")) == 2 {
        fmt.Println("Saving")
        response.Send("SAVING", 0)
        playHolder, err := response.RecvBytes(0)
        err = bson.Unmarshal(playHolder, &play)
        if err != nil {
          panic(err)
        }
        //fmt.Println("\033[38:2:200:0:200m",play,"\033[0m")
        savePfile(play)
        response.Send("SAVED", 0)
      }
    }else if strings.Contains(request, "+==LOGOUT") {
        playerToLogOut := strings.Split(request, "+==LOGOUT")[0]
        for i := 0;i < len(playerList);i++ {
          if strings.ToLower(playerList[i].Name) == strings.ToLower(playerToLogOut) {
            fmt.Println(playerList[i].Name+ "LOGGED OUT")
            response.Send("LOGGED "+playerList[i].Name+" OUT", 0)
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
    }else if strings.Contains(request, "=+=") {
      broadcastContainer = getBroadcasts()
      out := ""
      outVal := ""
      row := 0
//      fmt.Println(broadcastContainer)
      for i := 0;i < len(broadcastContainer);i++ {
        broadcastContainer[i].Payload.ID = i
        outVal = AssembleBroadside(broadcastContainer[i], broadcastContainer[i].Payload.Row, broadcastContainer[i].Payload.Col)
        row += 4
        out += outVal
      }
      _, err := response.Send(out, 0)
      if err != nil {
        panic(err)
      }
      }else if strings.Contains(request, "--UPSERT--") {
        _, err = response.Send("OKTOSEND", 0)
        if err != nil {
          panic(err)
        }
        socBytes, err := response.RecvBytes(0)
        if err != nil {
          panic(err)
        }
        var broac Broadcast
        err = json.Unmarshal(socBytes, &broac)
        if err != nil {
          panic(err)
        }
      //    fmt.Println("THIS SHOULD NOT BE ZERO")
      //      fmt.Println(broadcastContainer)
          insertBroadcast(broac)
          _, err = response.Send("DONE", 0)
    }else if strings.Contains(request, "--+--") {
      _, err = response.Send("OKTOSEND", 0)
      if err != nil {
        panic(err)
      }
      socBytes, err := response.RecvBytes(0)
      if err != nil {
        panic(err)
      }

      err = json.Unmarshal(socBytes, &broadcastContainer)
      if err != nil {
        panic(err)
      }
  //    fmt.Println("THIS SHOULD NOT BE ZERO")
//      fmt.Println(broadcastContainer)
  //    insertBroadcast(broadcastContainer)
      out := ""
      BROAD:
      for count := 0;count < len(broadcastContainer);count++ {
        for row := 0;row <= 20;row += 4 {
          for col := 53;col <= 143;col += 30 {
            if count >= len(broadcastContainer) {
              break BROAD
            }
            broadcastContainer[count].Payload.ID = count
            broad := AssembleBroadside(broadcastContainer[count], row, col)
            broadcastContainer[count].Payload.Row = row
            broadcastContainer[count].Payload.Col = col
            out += broad
            count++
          }
        }

      }

      _, err = response.Send(out, 0)
      if err != nil {
        panic(err)
      }
      }else if strings.HasPrefix(request, "--SELECT:") {
        fmt.Println("GOT SELECTOR")
        fmt.Println("Clearing selection")
        for i := 0;i < len(broadcastContainer);i++ {
          broadcastContainer[i].Payload.Selected = false
        }
        numSelected, err := strconv.Atoi(strings.Split(request, "--SELECT:")[1])
        if err != nil {
          match := strings.Split(request, "--SELECT:")[1]
          fmt.Println("Non-integer index! Trying a fuzzy match!")
          for i := 0;i < len(broadcastContainer);i++ {
            if strings.Contains(broadcastContainer[i].Payload.Message, match) {
              broadcastContainer[i].Payload.Selected = true
            }else if strings.Contains(broadcastContainer[i].Ref, match) {
              broadcastContainer[i].Payload.Selected = true
            }
          }
          numSelected = -1
        }
        //fmt.Println(string(request))
        for i := 0;i < len(broadcastContainer);i++ {
          if i == numSelected {
            broadcastContainer[i].Payload.Selected = true
          }
        }
        broadBytes, err := json.Marshal(broadcastContainer)
        if err != nil {
          panic(err)
        }
        _, err = response.SendBytes(broadBytes, 0)
        if err != nil {
          panic(err)
        }
        fmt.Println("Send ok")
      }else {
        broadcastContainer = getBroadcasts()
  //    in <- request
      _, err := response.Send("INVALID REQUEST", 0)
  //    fmt.Println("\033[38:2:150:0:150m"+request+"\033[0m")
      if err != nil {
        panic(err)
      }
    }
  }

}

func AssembleBroadside(broadside Broadcast, row int, col int) (string) {
	var cel string
	colString := strconv.Itoa(col)
	inWord := broadside.Payload.Message
	wor := ""
	word := ""
	words := ""
	if len(inWord) > 68 {
		return "DONE COMPOSTING"
	}
	if len(inWord) > 28 && len(inWord) > 54 {
		wor += inWord[:28]
		word += inWord[28:54]
		words += inWord[54:]
		for i := len(words); i <= 28; i++ {
			words += " "
		}
	}
	if len(inWord) > 28 && len(inWord) < 54 {
		wor += inWord[:28]
		word += inWord[28:]
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "

	}
	if len(inWord) <= 28 {
		wor = "                            "
		word += ""
		word += inWord
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "
	}

  	row++
  	if broadside.Payload.Selected {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m ",broadside.Ref, wor[len(broadside.Ref):], "\033[48;2;200;25;150m \033[0m")
  	}else {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m",broadside.Ref, wor[len(broadside.Ref):], "\033[48;2;20;255;50m \033[0m")
  	}

  	row++
  	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", word, "\033[48;2;20;255;50m \033[0m")
  	row++
  	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", words, "\033[48;2;20;255;50m \033[0m")
  	row++
  	if broadside.Payload.Game == "" {
  		if broadside.Payload.Selected {
  			broadside.Payload.Game = "SELECTED"
  		} else {
  			broadside.Payload.Game = "snowcrash"
  		}
  	}

    numString := strconv.Itoa(broadside.Payload.ID)
  	namePlate := "                            "[len(broadside.Payload.Name+numString):]
  	if broadside.Payload.Selected {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m @"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;200;25;50m \033[0m")
  	}else {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m@"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;20;255;50m \033[0m")

  	}
  	broadRow := 0
  	if broadside.Payload.Selected && len(strings.Split(broadside.Payload.BigMessage, "\n")) > 1 {
  		bigSplit := strings.Split(broadside.Payload.BigMessage, "\n")
  		for i := 0;i < len(bigSplit);i++ {
  			cel += fmt.Sprint("\033["+strconv.Itoa(25+broadRow)+";53H\033[48:2:200:0:200m \033[0m"+bigSplit[broadRow]+"\033[48:2:200:0:200m \033[0m")
  			broadRow++
  		}
  	}
	return cel
	//	fmt.Println(cel)
}
func main() {
    broadcast := make([]Broadcast, 1)
    broadcastLine := make(chan Broadcast)
    vineOn := make(chan bool)
    in := make(chan string)
    var populated []Space
    var wizinit bool
    wizinit = false
    playerList := make(chan string)
    populated = PopulateAreas()
    grapeVine(broadcast, playerList, vineOn, broadcastLine)
    if len(os.Args) > 1 {
      if len(os.Args) == 2 && os.Args[1] == "--wizinit" {
        dorp := InitPlayer("dorp", "norp")
        descString := "The absence of light is blinding.\nThree large telephone poles illuminate a small square."

        InitZoneSpaces("5-15", "Midgaard", descString)
        populated = PopulateAreas()
        addPfile(dorp, "norp")
        dorp.CurrentRoom = populated[0]
        savePfile(dorp)
        fmt.Println("dorp/norp")
        wizinit = true
      }
    }
    if !wizinit {
      go loopInput(populated, broadcast, in, playerList, vineOn, broadcastLine)
    }
    for {
      if !wizinit {
        <- in
      }
        if len(os.Args) == 2 && os.Args[1] == "--wizinit" {
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

  			createMobiles("Noodles")
        respond := fmt.Sprint("\033[38:2:0:250:0mInitialized "+strconv.Itoa(len(populated))+" rooms\033[0m")
        fmt.Printf(respond)

        fmt.Println("\033[38:2:0:250:0mAll tests passed and world has been initialzed\n\033[0mYou may now start by running ./sncn-core with no arguments and logging in from a client\nUser:dorp\nPass:norp.")
        os.Exit(1)
      }
    }
}

func UIDMaker() string {
	hostname := "localhost"
	username := "username"
	//Inspired by 'Una (unascribed)'s bikeshed
	rand.Seed(int64(time.Now().Nanosecond()))
	adjectives := []string{"Accidental", "Allocated", "Asymptotic", "Background", "Binary",
		"Bit", "Blast", "Blocked", "Bronze", "Captured", "Classic",
		"Compact", "Compressed", "Concatenated", "Conventional",
		"Cryptographic", "Decimal", "Decompressed", "Deflated",
		"Defragmented", "Dirty", "Distinguished", "Dozenal", "Elegant",
		"Encrypted", "Ender", "Enhanced", "Escaped", "Euclidean",
		"Expanded", "Expansive", "Explosive", "Extended", "Extreme",
		"Floppy", "Foreground", "Fragmented", "Garbage", "Giga", "Gold",
		"Hard", "Helical", "Hexadecimal", "Higher", "Infinite", "Inflated",
		"Intentional", "Interlaced", "Kilo", "Legacy", "Lower", "Magical",
		"Mapped", "Mega", "Nonlinear", "Noodle", "Null", "Obvious", "Paged",
		"Parity", "Platinum", "Primary", "Progressive", "Prompt",
		"Protected", "Quick", "Real", "Recursive", "Replica", "Resident",
		"Retried", "Root", "Secure", "Silver", "SolidState", "Super",
		"Swap", "Switched", "Synergistic", "Tera", "Terminated", "Ternary",
		"Traditional", "Unlimited", "Unreal", "Upper", "Userspace",
		"Vector", "Virtual", "Web", "WoodGrain", "Written", "Zipped"}
	nouns := []string{"AGP", "Algorithm", "Apparatus", "Array", "Bot", "Bus", "Capacitor",
		"Card", "Chip", "Collection", "Command", "Connection", "Cookie",
		"DLC", "DMA", "Daemon", "Data", "Database", "Density", "Desktop",
		"Device", "Directory", "Disk", "Dongle", "Executable", "Expansion",
		"Folder", "Glue", "Gremlin", "IRQ", "ISA", "Instruction",
		"Interface", "Job", "Key", "List", "MBR", "Map", "Modem", "Monster",
		"Numeral", "PCI", "Paradigm", "Plant", "Port", "Process",
		"Protocol", "Registry", "Repository", "Rights", "Scanline", "Set",
		"Slot", "Smoke", "Sweeper", "TSR", "Table", "Task", "Thread",
		"Tracker", "USB", "Vector", "Window"}
	uniquefier := ""

	uniqe := ""
	for i := 0; i < 2; i++ {
		uniq := rand.Intn(15)
		if uniq >= 10 {
			switch uniq {
			case 10:
				uniqe = "A"
			case 11:
				uniqe = "B"
			case 12:
				uniqe = "C"
			case 13:
				uniqe = "D"
			case 14:
				uniqe = "E"
			case 15:
				uniqe = "F"
			}

			uniquefier += uniqe
		} else {
			uniquefier += fmt.Sprint(uniq)
		}
	}
	ind := rand.Intn(len(adjectives))
	indie := rand.Intn(len(adjectives))
	if indie == ind {
		indie = rand.Intn(len(adjectives))
	}
	thedog := rand.Intn(len(nouns))
	uniqueFied := fmt.Sprint(uniquefier, adjectives[ind], adjectives[indie], nouns[thedog])

	//fmt.Println(uniqueFied)

	UID := fmt.Sprint(uniqueFied, hostname, username)
	return UID
}
