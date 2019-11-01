package main

import (
  "fmt"
  zmq "github.com/pebbe/zmq4"

)

func main() {
  ctx, err := zmq.NewContext()
  if err != nil {
    panic(err)
  }

  inSocket, err := ctx.NewSocket(zmq.REP)
  if err != nil {
    panic(err)
  }
  outSocket, err := ctx.NewSocket(zmq.REQ)
  if err != nil {
    panic(err)
  }
  inSocket.Bind("tcp://127.0.0.1:4000")
  outSocket.Connect("tcp://127.0.0.1:4000")

  outSocket.Send("Noot!", 0)
  reply, err := inSocket.Recv(0)
  if err != nil {
    panic(err)
  }

  fmt.Println(reply)

  fmt.Println("Let's fill this space with the core functionality")
}
