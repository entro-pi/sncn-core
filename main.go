package main

import (
  "time"
  "strconv"
  "strings"
  "fmt"
  "os"
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
  response, err := zmq.NewSocket(zmq.REP)
  if err != nil {
    panic(err)
  }
  defer response.Close()
  //Preferred way to connect
  //hostname := "tcp://snowcrashnetwork.vineyard.haus:7778"
  hostname := "tcp://*:7777"
//  clientname := "tcp://192.168.122.1:7777"
  err = response.Bind(hostname)
  if err != nil {
    panic(err)
  }
//  err = login.Connect(hostname)

  for {
    fmt.Println("IN LOOP")
    request, err := response.Recv(0)
    if err != nil {
      panic(err)
    }
    fmt.Println(string(request))
    if strings.Split(string(request), ":")[0] == "REQUESTPUBKEY" {

        _, err = response.Send(servepubKey, 0)
        in <- request
        if err != nil {
          panic(err)
        }
    //    resp, err := response.Recv(0)
    //    if err != nil {
    //      panic(err)
    //    }
    //    in <- string(resp)
    }else if string(request) == "shutdown" {
      os.Exit(1)
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
    go givePubKey("Doot", in)
    for {
      value := <-in
      fmt.Println("\033[38:2:255:0:0m"+value+"\033[0m")
    }
  fmt.Println("Let's fill this space with the core functionality")
}
