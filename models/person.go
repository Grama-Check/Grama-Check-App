package models

type Person struct {
	NIC     string `json:"nic"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type IDChecked struct {
	NIC    string `json:"nic"`
	Exists bool   `json:"exists"`
}

type AddressChecked struct {
	NIC    string `json:"nic"`
	Exists bool   `json:"exists"`
}

type PoliceCheck struct {
	NIC   string `json:"nic"`
	Clear string `json:"clear"`
}
