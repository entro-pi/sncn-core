package main

import (
  "time"
  "math/rand"
  "strconv"
  "strings"
  "fmt"
  "os"
  "bufio"
  "io/ioutil"
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  zmq "github.com/pebbe/zmq4"
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
//  go client(broadcast, string(clientid), string(secret), grapevine, playerList, vineOn, broadcastLine)
}

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*32+100)
  }
  return newVal
}
func onlineHash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*24+240)
  }
  return newVal
}
func updatePlayerSlain(hash string) {
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
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
  _, err = collection.UpdateOne(context.Background(), options.Update().SetUpsert(false), bson.M{"$set":bson.M{"name":player.Name,"title":player.Title,"inventory":player.Inventory, "equipped":player.Equipped,
						"coreboard": player.CoreBoard,"currentroom":player.CurrentRoom,"slain":player.Slain, "hoarded":player.Hoarded, "str": player.Str, "int": player.Int, "dex": player.Dex, "wis": player.Wis, "con":player.Con, "cha":player.Cha, "classes": player.Classes }})

  if err != nil {
    panic(err)
  }







}
func lookupPlayerByHash(playerHash string) Player {
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
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

  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":playerHash}})
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
func lookupPlayer(name string, pass string) Player {
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  passCred := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+passCred+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
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

func onlineTransaction(advert *Broadcast, customer Player, allItems []Object) (Player, string) {
  output := ""
  if advert.Payload.Transaction.Sold {
    return customer, "SOLD OUT"
  }
  if len(advert.Payload.Store.Inventory) > 0 {
      for i := 0;i < len(advert.Payload.Store.Inventory);i++ {
        hasSpace := false
        slot := 0
        vnum := advert.Payload.Store.Inventory[i].Item.Vnum
        hash := advert.Payload.Store.Inventory[i].ItemHash
        price := advert.Payload.Store.Inventory[i].Price
        customerCash := customer.BankAccount.Amount
        isSold := advert.Payload.Store.Inventory[i].Sold
        fmt.Println("VNUM",vnum,"HASH",hash,"PRICE",price,"CUSTOMERCASH",customerCash,"ISSOLD",isSold)
        for c := len(customer.Inventory) - 1;c > 0;c-- {
          if customer.Inventory[c].Item.Name == "nothing" {
            hasSpace = true
            slot = c
          }
        }
        if customerCash >= price && hasSpace {
            customer.BankAccount.Amount -= price
        }
        if hash == onlineHash(allItems[vnum].LongName) {
          customer.Inventory[slot].Item = allItems[vnum]
          customer.Inventory[slot].Number++
          advert.Payload.Store.Inventory[i].Sold = true
          fmt.Println("\033[38:2:0:200:0mTransaction approved.\033[0m")
          output += fmt.Sprintln("\033[38:2:0:200:0mTransaction approved.\033[0m")

          return customer, output
        }
      }
  }else if !advert.Payload.Transaction.Sold {
      hasSpace := false
      hasCash := false
      slot := 0
      vnum := advert.Payload.Transaction.Item.Vnum
      hash := advert.Payload.Transaction.ItemHash
      price := advert.Payload.Transaction.Price
      customerCash := customer.BankAccount.Amount

      fmt.Println("VNUM",vnum,"HASH",hash,"PRICE",price,"CUSTOMERCASH",customerCash)
      for c := len(customer.Inventory) - 1;c > 0;c-- {
        if customer.Inventory[c].Item.Name == "nothing" || customer.Inventory[c].Item.Name == "" || customer.Inventory[c].Item.Name == advert.Payload.Transaction.Item.Name{
          hasSpace = true
          slot = c
        }
      }
      if customerCash >= price && hasSpace {
          customer.BankAccount.Amount -= price
          hasCash = true
      }
      if hasSpace && hasCash {
        customer.Inventory[slot].Item = allItems[vnum]
        customer.Inventory[slot].Number++
        fmt.Println("\033[38:2:0:200:0mTransaction approved.\033[0m")
        output += fmt.Sprintln("\033[38:2:0:200:0mTransaction approved.\033[0m")
        return customer, output

      }
  }else {
    output += fmt.Sprint("Looks like you missed out on the sale!")
    output += fmt.Sprint("That is sold out!")
  }
  output += fmt.Sprintln("\033[38:2:200:0:0mTransaction declined.\033[0m")
  fmt.Println("\033[38:2:200:0:0mTransaction declined.\033[0m")
  return customer, output
}

func onlineButlerTransaction(advert Broadcast, customer Butler) Object {
  //let's get players up first
  var blank Object
  return blank
}


func getPlayers() []Player {
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }

  collection := client.Database("pfiles").Collection("Players")
  curs, err := collection.Find(context.Background(), bson.M{})
  if err != nil {
    panic(err)
  }
  var container []Player
  for curs.Next(context.Background()) {
    var play Player
    err = curs.Decode(&play)
    if err != nil {
      panic(err)
    }
    container = append(container, play)
  }

  return container

}

func getGrapes() []Broadcast {
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
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
  userFile, err := os.Open("weaselcreds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@cloud-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  update := bson.M{"event":bcast.Event,"ref":bcast.Ref,"payload":bson.M{"channel":bcast.Payload.Channel,"id":bcast.Payload.ID, "message":bcast.Payload.Message,"game":bcast.Payload.Game,"name":bcast.Payload.Name,"transaction":bcast.Payload.Transaction}}
  collection := client.Database("broadcasts").Collection("snowcrash")
  _, err = collection.InsertOne(context.Background(), update)
  if err != nil {
    panic(err)
  }
  fmt.Println("Upserted the broadcast")
  return bcast

}

func loopInput(populated []Space, broadcast []Broadcast, in chan string, players chan string, vineOn chan bool, broadcastLine chan Broadcast, allItems []Object) {
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
  request := ""
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
        //time.Sleep(15*time.Second)
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
      request = ""
    }
    fmt.Println("IN LOOP")
    //request, err = response.Recv(0)
    //if err != nil {
    //  panic(err)
    //}


    var play Player
    select {
    case request = <- in:
      _, err = response.Recv(0)
      if err != nil {
        panic(err)
      }
    default:
      request, err = response.Recv(0)
      if err != nil {
        panic(err)
      }
    }
    if strings.Contains(request, ":"){
      fmt.Println(strings.Split(request, ":")[0]+"AUTHINFO")
    }else if strings.Contains(request, "+=+") {
      fmt.Println(strings.Split(request, "+=+")[0]+"CHATMSGS")
    }else {
      fmt.Println(request)
    }

    if strings.Contains(request, "\"event\": \"restart\",") {
      //time.Sleep(15*time.Second)
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
    if strings.Contains(request, "||SALE||") {
      fmt.Println("\033[38:2:150:0:150mSale going down!\033[0m")
      playerHash := strings.Split(request, "||SALE||")[0]
      ref := strings.Split(request, "||SALE||")[1]
      var advert Broadcast
      //TODO
      broadcasts := getBroadcasts()
      found := false
      GETADVERT:
      for i := 0;i < len(broadcasts);i++ {
        if broadcasts[i].Ref == ref {
          advert = broadcasts[i]
          fmt.Println("BROADCAST IS ",broadcasts[i].Payload)
          fmt.Println("ADVERT VNUM IS",advert.Payload.Transaction.Item.Name)
          found = true
          break GETADVERT
        }
      }
      if found {
        play = lookupPlayerByHash(playerHash)
        fmt.Println("ADVERT VNUM IS",advert.Payload.Transaction.Item.Name)
        play, output := onlineTransaction(&advert, play, allItems)
        if strings.Contains(output, "approved") {
          advert.Payload.Transaction.Sold = true
        }
        updateBroadcast(advert)
        _, err = response.Send(output, 0)
        _, err := response.Recv(0)
//        fmt.Println(result)
        playBytes, err := bson.Marshal(play)
        if err != nil {
          panic(err)
        }
        fmt.Println("GOT TO SEND PLAYER")
        _, err = response.SendBytes(playBytes, 0)
        fmt.Println("SENT PLAYER")
        if err != nil {
          panic(err)
        }
        response.Recv(0)
        request = ""
      }
    }
    if strings.Contains(request, ":-:") {
        userPass := strings.Split(request, ":-:")
        name, pass := userPass[0], userPass[1]
        play = InitPlayer(name, pass)
        play.CurrentRoom = populated[0]
        addPfile(play, pass)
//        savePfile(play)
        playBytes, err := bson.Marshal(play)
        _, err = response.SendBytes(playBytes, 0)
        if err != nil {
          panic(err)
        }
/*    }else if strings.Contains(request, ":=:") {
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
        players <- play.Name*/
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

    }else if strings.Contains(request, "+=+") {


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
      //  savePfile(play)
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
            stringID := strconv.Itoa(broadcastContainer[i].Payload.ID)
            if strings.Contains(broadcastContainer[i].Payload.Message, match) {
              broadcastContainer[i].Payload.Selected = true
            }else if strings.Contains(stringID, match) {
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
    payloadString := strconv.Itoa(broadside.Payload.ID)
  	row++
  	if broadside.Payload.Selected {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m ",payloadString, wor[len(payloadString):], "\033[48;2;200;25;150m \033[0m")
  	}else {
  		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m",payloadString, wor[len(payloadString):], "\033[48;2;20;255;50m \033[0m")
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
    //allItems := readItemsFromFile("dat/items/items.itm")
    playerSignIn, err := zmq.NewSocket(zmq.REP)
    if err != nil {
      panic(err)
    }
    battling, err := zmq.NewSocket(zmq.PUB)
    if err != nil {
      panic(err)
    }
    var players []Player
    playerSignIn.Bind("tcp://127.0.0.1:7776")
    battling.Bind("tcp://127.0.0.1:7777")
    for {
      incomingPlayer, err := playerSignIn.Recv(0)
      if err != nil {
        panic(err)
      }
        if strings.Contains(incomingPlayer, "ACCEPT") {
          playHash := strings.Split(incomingPlayer, ":")[1]
          player := lookupPlayerByHash(playHash)
          if player.Name != "" {
            players = append(players, player)
          }
        }
        if strings.Contains(incomingPlayer, "ATTACK") {
          splitNameAttacked := strings.Split(incomingPlayer, ";")
          playName, mobAttacked := splitNameAttacked[1], splitNameAttacked[2]
          for i := 0;i < len(players);i++ {
            if players[i].Name == playName {
              players[i].Battling = true
              play := players[i]
              for c := 0;c < len(play.Fights.Oppose);c++ {
                if play.Fights.Oppose[c].Rezz > 0 && play.Fights.Oppose[c].Name == mobAttacked {
                  //FIGHT!
                  act := PvEResolve(players[i])
                  _, err := battling.Send(play.Name+":"+act.DamMsg+"="+strconv.Itoa(act.Damage), 0)
                  if err != nil {
                    panic(err)
                  }
                }else if play.Fights.Oppose[c].Rezz <= 0 {
                  //already dead!
                  players[i].Battling = false
                }else {
                  //not here!
                  players[i].Battling = false
                }

              }
            }
          }
        }
      //This is just because we can't properly trigger it yet


    }

}

func PvEResolve(play Player) Action {
  var act Action
  mob := play.BattlingMob
  act.Affects = play
  act.From = mob
  Attack := rand.Intn(act.From.Attack)
  Defend := rand.Intn(act.Affects.Defend) - Attack
  if Defend < 0 {
    act.Damage = Defend * -1 + rand.Intn(act.From.Attack)
    act.DamMsg = mob.Name+" swings widly and lands \033[38:2:200:0:0m"+ strconv.Itoa(act.Damage)+ "\033[0m points of damage!"
  }else {
    act.Damage = 0
    act.DamMsg = mob.Name +" \033[38:2:0:200:0m fails to do any damage\033[0m!"
  }
  return act
}

func playerListener(incoming chan string, outgoing chan Player) {
  for {
    select {
    case inValue := <-incoming:
      players := getPlayers()
      inSplit := strings.Split(inValue, ":")
      if inSplit[0] == "ACCEPT" {
          for i := 0;i < len(players);i++ {
            if players[i].PlayerHash == inSplit[1] {
              //Success
              outgoing <- players[i]
            }
          }
      }
    }
  }

}
func playerBattle(play Player) {
  out, err := zmq.NewSocket(zmq.PUB)
  if err != nil {
    panic(err)
  }
  err = out.Bind("tcp://127.0.0.1:7777")
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
