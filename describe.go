package main

import (
  "fmt"
  "strconv"
)
func clearDirty() {
  for i := 0;i < 255;i++ {
    fmt.Println("")
  }
}

func clearCmd() {
		fmt.Print(cmdPos+"                                                                                                                                                                                   ")
		fmt.Print("\033[52;0H                                                                                                                                                                                   ")
		fmt.Print("\033[53;0H                                                                                                                                                                                   ")
		fmt.Print("\033[54;0H                                                                                                                                                                                   ")
		fmt.Print("\033[55;0H                                                                                                                                                                                   ")
		fmt.Print("\033[56;0H                                                                                                                                                                                   ")
		fmt.Print(cmdPos)
}

func drawDig(digFrame [][]int, zonePos []int) {
	for i := 0;i < len(digFrame);i++ {
		for c := 0;c < len(digFrame[i]);c++ {
				prn := ""
				val := fmt.Sprint(digFrame[i][c])
				if i == zonePos[0] && c == zonePos[1] {
					prn = "8"
				}
				if prn == "8" {
					fmt.Printf("\033[38:2:150:10:50m"+val+"\033[0m")
				}else if val == "1" || val == "8" {
					val = "1"
					fmt.Printf("\033[38:2:50:10:50m"+val+"\033[0m")
				}else {
						fmt.Printf(val)
				}
		}
		fmt.Println("")
	}
}

func DescribePlayer(play Player) {

  ratio := ""
  count := 18
  for   rezz := 0;rezz < play.Rezz;rezz++ {

    ratio += "\033["+strconv.Itoa(count+30)+";25H\033[48:2:175:50:50m \033[0m\n"
    count--
  }
  for count > 0 {
      ratio += "\033["+strconv.Itoa(count+30)+";25H\033[48:2:15:50:50m \033[0m\n"

    count--
  }

  ratio += "\033[31;24H+++\n"
  ratio += "\033[49;24H+++"
  hp := ratio
  count = 18
  ratio = ""
  for tech := 0;tech < play.Tech;tech++ {
    ratio += "\033["+strconv.Itoa(count+30)+";31H\033[48:2:75:150:50m \033[0m\n"
    count--
  }
  for count > 0 {
      ratio += "\033["+strconv.Itoa(count+30)+";31H\033[48:2:15:50:50m \033[0m\n"
      count--
  }
  ratio += "\033[31;30H===\n"
  ratio += "\033[49;30H==="
  techShow := ratio

  fmt.Print(techShow)
  fmt.Print(hp)
	fmt.Printf("\033[40;0H")
	fmt.Println("======================")
	fmt.Println("\033[38:2:0:200:0mStrength     :\033[0m", play.Str)
	fmt.Println("\033[38:2:0:200:0mIntelligence :\033[0m", play.Int)
	fmt.Println("\033[38:2:0:200:0mDexterity    :\033[0m", play.Dex)
	fmt.Println("\033[38:2:0:200:0mWisdom       :\033[0m", play.Wis)
	fmt.Println("\033[38:2:0:200:0mConstitution :\033[0m", play.Con)
	fmt.Println("\033[38:2:0:200:0mCharisma     :\033[0m", play.Cha)
	fmt.Println("======================")
}
