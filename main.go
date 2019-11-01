package main

import (
  "time"
  "strconv"
  "strings"
  "fmt"
  "github.com/SolarLune/dngn"
  zmq "github.com/pebbe/zmq4"

)


type Descriptions struct {
	BATTLESPAM int
	ROOMDESC int
	PLAYERDESC int
	ROOMTITLE int
}
type Chat struct {
	User Player
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

)

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*16+25)
  }
  return newVal
}

func givePubKey(servepubKey string, in chan string) {
  fmt.Println("Core login procedure started")
  login, err := zmq.NewSocket(zmq.PUSH)
  if err != nil {
    panic(err)
  }
  defer login.Close()
  response, err := zmq.NewSocket(zmq.PULL)
  if err != nil {
    panic(err)
  }
  defer response.Close()
  //Preferred way to connect
  //hostname := "tcp://snowcrashnetwork.vineyard.haus:4000"
  hostname := "tcp://*:4001"
//  clientname := "tcp://192.168.122.1:4001"
  err = response.Bind(hostname)
//  err = login.Connect(hostname)

  for {
    fmt.Println("IN LOOP")
    request, err := response.Recv(0)
    if err != nil {
      panic(err)
    }
    fmt.Println(string(request))
    if strings.Split(string(request), ":")[0] == "REQUESTPUBKEY" {
        err := login.Connect("tcp://"+strings.Split(string(request), ":")[1]+":4001")
        if err != nil {
          panic(err)
        }
        _, err = login.Send(servepubKey, 0)
        in <- request
        if err != nil {
          panic(err)
        }
        resp, err := response.Recv(0)
        if err != nil {
          panic(err)
        }
        in <- string(resp)
    }else {

      _, err := login.Send("INVALID REQUEST", 0)
      if err != nil {
        panic(err)
      }
    }
  }

}

func main() {
  in := make(chan string)

  client, err := zmq.NewSocket(zmq.PULL)
  if err != nil {
    panic(err)
  }
  defer client.Close()

//  server.SetSockOpt(zmq.ZMQ_CURVE_SERVER)
  server, err := zmq.NewSocket(zmq.PUSH)
  if err != nil {
    panic(err)
  }
  defer server.Close()
  if zmq.HasCurve() {

    zmq.AuthSetVerbose(true)
    zmq.AuthStart()
    servekey, servesec, err := zmq.NewCurveKeypair()
    //have this run as it's own thread
    go givePubKey(servekey, in)
    clientDat := <-in
    clientString := strings.Split(clientDat, ":")
    _, IPAddress := clientString[0], clientString[1]


    zmq.AuthAllow("snowcrash.network", IPAddress+"/8")

//    go givePubKey(servekey)
    zmq.AuthCurveAdd("snowcrash.network", zmq.CURVE_ALLOW_ANY )
    err = server.ServerAuthCurve("snowcrash.network", servesec)
    server.Bind("tcp://*:4000")

    if err != nil {
      panic(err)
    }
    fmt.Println("Connecting...")
    //channel in!
    client.Connect("tcp://"+IPAddress+":4000")
    time.Sleep(100*time.Millisecond)
    fmt.Println("Sending...")
    _, err = server.Send("Curve security status: True", 0)
    if err != nil {
      panic(err)
    }
  }else {
    server.Bind("tcp://127.0.0.1:4000")
    time.Sleep(100*time.Millisecond)
    server.Send("Curve security status: False", 0)
  }
  for {
    command, err := client.Recv(0)
    if err != nil {
      panic(err)
    }
    fmt.Println("\033[38:2:255:0:0mINPUT WAS"+command+"\033[0m")
      if command == "shutdown" {

      }
  }
  zmq.AuthStop()


  fmt.Println("Let's fill this space with the core functionality")
}
