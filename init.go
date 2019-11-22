package main

import (

  "context"
  "time"
	"strconv"
	"strings"
  "os"
  "bufio"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func initDigRoom(digFrame [][]int, zoneVnums string, zoneName string, play Player, vnum int) (Space, int) {
	var dg Space
	dg.Vnums = zoneVnums
	dg.Zone = zoneName
	dg.ZonePos = make([]int, 2)
	dg.ZoneMap = digFrame
	//todo directions
	vnum += 1
	dg.Vnum = vnum
	dg.Altered = true
	dg.Desc = "Nothing but some cosmic rays"
	for len(strings.Split(dg.Desc, "\n")) < 8 {
		dg.Desc += "\n"
	}
	return dg, vnum
}



func InitPlayer(name string, pass string) Player {
	var play Player

  var class Class
	play.Name = name
  play.PlayerHash = hash(name+pass)
  play.Title = "The Unknown"
  play.Classes = append(play.Classes, class)
  play.Classes[0].Level = 1
  play.Classes[0].Name = "wildling"
  var rip Skill
  rip.DamType = "slash"
  rip.Name = "overcharge"
  rip.Level = 1
  rip.Usage = 'e'
  play.Classes[0].Skills = append(play.Classes[0].Skills, rip)
  var blast Spell
  blast.TechUsage = 2
  blast.Name = "blast"
  blast.Level = 1
  blast.Dam = 3
  blast.Consumed = false
  play.Classes[0].Spells = append(play.Classes[0].Spells, blast)

	play.Inventory = make([]InventoryItem, 20, 20)
  play.Equipped = make([]EquipmentItem, 20, 20)
  play.Rezz = 17
  play.MaxRezz = play.Rezz
  play.Tech = 17

	play.Str = 1
	play.Int = 1
	play.Dex = 1
	play.Wis = 1
	play.Con = 1
	play.Cha = 1
  play.Channels = append(play.Channels, "testing")
  savePfile(play)
	return play

}

func InitZoneSpaces(SpaceRange string, zoneName string, desc string) {
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
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("zones").Collection("Spaces")
	vnums := strings.Split(SpaceRange, "-")
	vnumStart, err := strconv.Atoi(vnums[0])
	if err != nil {
		panic(err)
	}

	vnumEnd, err := strconv.Atoi(vnums[1])
	if err != nil {
		panic(err)
	}
	for i := vnumStart;i < vnumEnd;i++ {
		var mobiles []int
		var items []int
		mobiles = append(mobiles, 0)
		items = append(items, 0)
		_, err = collection.InsertOne(context.Background(), bson.M{"vnums":SpaceRange,"zone":zoneName,"vnum":i, "desc":desc,
							"mobiles": mobiles, "items": items })
	}
	if err != nil {
		panic(err)
	}
}
