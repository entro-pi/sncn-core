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
    in <- request
    fmt.Println(string(request))
    if strings.Split(string(request), ":")[0] == "REQUESTPUBKEY" {

        _, err = response.Send(servepubKey, 0)
//        in <- request
        if err != nil {
          panic(err)
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


    clientkey, _, err := zmq.NewCurveKeypair()
    if err != nil {
      panic(err)
    }
    go loopInput(clientkey, in)
    for {
      value := <-in
      if value == "init world" {
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
  			play = InitPlayer("FSM")
  			addPfile(play)
  			createMobiles("Noodles")
        respond := fmt.Sprint("\033[38:2:0:250:0mInitialized "+strconv.Itoa(len(populated))+" rooms\033[0m")
        fmt.Printf(respond)

        fmt.Println("\033[38:2:0:250:0mAll tests passed and world has been initialzed\n\033[0mYou may now start with --login.")

      }else if value == "shutdown" {

        fmt.Println("\033[38:2:255:0:0mGOT "+value+" SIGNAL\033[0m")
        os.Exit(1)
      }
      fmt.Println("\033[38:2:0:150:0m"+value+"\033[0m")
    }
  fmt.Println("Let's fill this space with the core functionality")
}
