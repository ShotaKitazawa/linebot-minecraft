package domain

type LineClientConfig struct {
	GroupID       string
	ChannelSecret string
	ChannelToken  string
}

type Domain struct {
	//ログインしてるユーザ
	LoginUsers []User
	//whitelistなユーザ名
	WhitelistUsernames []string
}

type User struct {
	Name     string
	XpLevel  string
	Position Position
	Biome    string // Minecraft 1.16~
}

type Position struct {
	X float32
	Y float32
	Z float32
}
