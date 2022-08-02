package models

type Person struct {
	UID     int    `json:uid`
	ID      string `json:id`
	Address string `json:address`
}

type IDChecked struct {
	UID    int  `json:uid`
	Exists bool `json:exists`
}
