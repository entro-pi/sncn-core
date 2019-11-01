package main

import (
  "time"
  "strconv"
  "strings"
  "fmt"
  zmq "github.com/pebbe/zmq4"

)

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*16+25)
  }
  return newVal
}

func givePubKey(servepubKey string) {
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
    request, err := response.Recv(0)
    if err != nil {
      panic(err)
    }
    if strings.Split(string(request), ":")[0] == "REQUESTPUBKEY" {
        err := login.Connect("tcp://"+strings.Split(string(request), ":")[1]+":4001")
        if err != nil {
          panic(err)
        }
        _, err = login.Send(servepubKey, 0)

        if err != nil {
          panic(err)
        }
    }else {

      _, err := login.Send("INVALID REQUEST", 0)
      if err != nil {
        panic(err)
      }
    }
  }

}

func main() {


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
  for {
    _, _, err := zmq.NewCurveKeypair()
    if err != nil {
      panic(err)
    }
    givePubKey(string("dummykey"))
  }
  if zmq.HasCurve() {
    zmq.AuthSetVerbose(true)
    zmq.AuthStart()
    zmq.AuthAllow("snowcrash.network", "127.0.0.1/8")

    clientkey, clientseckey, err := zmq.NewCurveKeypair()
    servekey, servesec, err := zmq.NewCurveKeypair()
    //have this run as it's own thread
//    go givePubKey(servekey)

    zmq.AuthCurveAdd("snowcrash.network", clientkey )
    err = client.ClientAuthCurve(servekey, clientkey, clientseckey)
    err = server.ServerAuthCurve("snowcrash.network", servesec)
    server.Bind("tcp://*:4000")

    if err != nil {
      panic(err)
    }
    fmt.Println("Connecting...")
    client.Connect("tcp://127.0.0.1:4000")
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
  reply, err := client.Recv(0)
  if err != nil {
    panic(err)
  }

  fmt.Println(reply)
  zmq.AuthStop()

  fmt.Println("Let's fill this space with the core functionality")
}
