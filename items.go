package main

import (
  "os"
  "fmt"
  "bufio"
  "strconv"
)

func readItemsFromFile(filePath string) []Object {
  file, err := os.Open(filePath)
  if err != nil {
    panic(err)
  }
  fmt.Println("LOADING ITEMS")
  var objHolder []Object
  scanner := bufio.NewScanner(file)
  var obj Object

  for scanner.Scan() {

    if scanner.Text() == "VNUM" {
      scanner.Scan()
      //fmt.Println("VNUM")
      //fmt.Println(scanner.Text())
      obj.Vnum, err = strconv.Atoi(scanner.Text())
      if err != nil {
        panic(err)
      }
    }
    if scanner.Text() == "NAME" {
      scanner.Scan()
      //fmt.Println("NAME")
      //fmt.Println(scanner.Text())
      obj.Name = scanner.Text()
    }
    if scanner.Text() == "LONGNAME" {
      scanner.Scan()
      //fmt.Println("LONGNAME")
      //fmt.Println(scanner.Text())
      obj.LongName = scanner.Text()
    }
    if scanner.Text() == "ZONE" {
      scanner.Scan()
      //fmt.Println("ZONE")
      //fmt.Println(scanner.Text())
      obj.Zone = scanner.Text()
    }
    if scanner.Text() == "VALUE" {
      scanner.Scan()
      //fmt.Println("VALUE")
      //fmt.Println(scanner.Text())
      obj.Value, err = strconv.Atoi(scanner.Text())
      if err != nil {
        panic(err)
      }
    }
    if scanner.Text() == "OWNED" {
      scanner.Scan()
      //fmt.Println("OWNED")
      //fmt.Println(scanner.Text())
      if scanner.Text() == "true"{
        obj.Owned = true
      }else {
        obj.Owned = false
      }
    }
    if scanner.Text() == "SLOT" {
      scanner.Scan()
      //fmt.Println("SLOT")
      //fmt.Println(scanner.Text())
      obj.Slot, err = strconv.Atoi(scanner.Text())
      if err != nil {
        panic(err)
      }
      objHolder = append(objHolder, obj)

    }
  }
  fmt.Println("DONE LOADING ITEMS")
  return objHolder
}
