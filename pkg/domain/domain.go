package domain

import "github.com/line/line-bot-sdk-go/linebot"

var (
	MinecraftHostname = ``
)

type LineClientConfig struct {
	Client   *linebot.Client
	GroupIDs []string
}

type Entity struct {
	//ログインしたことのあるすべてのユーザ
	AllUsers []User
	//ログインしてるユーザ
	LoginUsers []User
	//ログインしてないユーザ
	LogoutUsers []User
	//whitelistなユーザ名
	WhitelistUsernames []string
}

type User struct {
	Name     string
	Health   float32
	XpLevel  uint
	Position Position
	Biome    string // Minecraft 1.16~
}

type Position struct {
	X float32
	Y float32
	Z float32
}
