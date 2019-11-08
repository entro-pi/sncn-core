package main

import (
	"github.com/SolarLune/dngn"
	"time"
)
type Class struct {
	Level float64
	Name string
	Skills []Skill
	Spells []Spell
}

type Spell struct {
	TechUsage int
	Name string
	Usage rune
	Dam int
	Level int
	Consumed bool
}
type Skill struct {
	Name string
	DamType string
	Level int
	Usage rune
}

type StatusPayload struct {
	Game string
	Players []string
}

type Status struct {
	Event string
	Ref string
	Payload StatusPayload
}


type SignOutPayload struct {
	Name string
	Game string
}

type SignOut struct {
	Event string
	Ref string
	Payload SignOutPayload
}

type SignIn struct {
	Event string
	Ref string
	Payload SignInPayload
}

type SignInPayload struct {
	Name string
	Game string
}

type BroadcastPayload struct {
  Channel string
  Message string
  Game string
  Name string
	Row int
	Col int
	Selected bool
}
type Broadcast struct {
    Event string
    Ref string
    Payload BroadcastPayload
}


type SendBPayload struct {
  Channel string
  Name string
  Message string
}

type SendBroadcast struct {
  Event string
  Ref string
  Payload SendBPayload
}

type Descriptions struct {
	BATTLESPAM int
	ROOMDESC int
	PLAYERDESC int
	ROOMTITLE int
}
type Chat struct {
	User string
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


type Object struct {
	Name string
	LongName string
	Vnum int
	Zone string
	Owner Player
	Value int
	X int
	Y int
	Owned bool
}
type Player struct {
	Name string
	Title string
	Inventory []int
	Equipment []int
	CoreBoard string
	PlainCoreBoard string
	CurrentRoom Space
	PlayerHash string
	Classes []Class
	Target string
	TargetLong string
	TarX int
	TarY int
	OldX int
	OldY int
	CPU string
	CoreShow bool
	Channels []string
	Battling bool
	Profile string
	Slain int
	Hoarded int

	MaxRezz int
	Rezz int
	Tech int
	Fights Fight
	Won int
	Found int

	Str int
	Int int
	Dex int
	Wis int
	Con int
	Cha int
}

type Fight struct {
	Oppose []Mobile
	Former []Player
	Treasure []Object
}
type Payload struct {
  Client_id string
  Client_secret string
  Supports []string
  Channels []string
  Version string
  User_agent string
  Unicode string
  Status string

}
type HeartPayload struct {
  Players []string
}

type Heartbeat struct {
  Event string
  Payload HeartPayload
}

type Authenticator struct {
  Event string
  Payload Payload
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
	X int
	Y int
	Char string
}
type GrapeMessPayload struct {
  Channel string
}

type GrapeMess struct {
  Event string
  Payload GrapeMessPayload
  Ref string
}
