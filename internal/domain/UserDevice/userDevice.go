package userDevice

type UserDevice struct {
	Id          string
	UserId      string
	DeviceName  string
	Platform    string
	LastSeen    string
	WebsocketId string
	IsActive    bool
}
