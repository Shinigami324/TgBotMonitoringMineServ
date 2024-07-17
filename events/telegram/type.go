package telegram

type PlayerData struct {
	Online int `json:"online"`
	Max    int `json:"max"`
}

type ServerData struct {
	Online bool       `json:"online"`
	Player PlayerData `json:"players"`
}
