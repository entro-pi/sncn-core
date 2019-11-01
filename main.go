package main

import (
  "time"
  "fmt"
  zmq "github.com/pebbe/zmq4"

)

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
  if zmq.HasCurve() {
    zmq.AuthSetVerbose(true)
    zmq.AuthStart()
    zmq.AuthAllow("snowcrash.network", "127.0.0.1/8")

    clientkey, clientseckey, err := zmq.NewCurveKeypair()
    servekey, servesec, err := zmq.NewCurveKeypair()

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
