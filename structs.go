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
type Bank struct {
	Clientele Client
	Owner string
}
//user and owner here will also be a playerhash
//on a UIDMaker() result rather than the pass
type Client struct {
		User string
		TotalAmount float64
		Accounts []float64
}
type Account struct {
	Owner string
	Amount float64
}

type BroadcastPayload struct {
  Channel string
  Message string
  Game string
  Name string
	Row int
	Col int
	Selected bool
	BigMessage string
	ID int
	Transaction OnlineTransaction
	Store OnlineStore
}

//Employer and owner will be a playerHash
type Butler struct {
	Employer string
	Funds Account
}
type OnlineStore struct {
	Owner string
	Float float64
	Inventory []OnlineTransaction
}

type OnlineTransaction struct {
	ItemHash string
	Item Object
	Sold bool
	To Account
	Price float64
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
	BigMessage string
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
	MobilesInRoom []Mobile
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
	Slot int
	Value int
	X int
	Y int
	Owned bool
}
type InventoryItem struct {
	Item Object
	Number int
}
type EquipmentItem struct {
	Item Object
}
type Player struct {
	Name string
	Title string
	Inventory []InventoryItem
	Equipped []EquipmentItem
	CoreBoard string
	PlainCoreBoard string
	CurrentRoom Space
	PlayerHash string
	Classes []Class
	Target string
	TargetLong string
	//ToBuy will have to either be an ItemHash
	//Or a vnum
	ToBuy int
	BankAccount Account
	TarX int
	TarY int
	OldX int
	OldY int
	CPU string
	CoreShow bool
	Channels []string
	Battling bool
	Profile string
	Session string

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
	Aggro bool
	Align int
	Vnum int
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
