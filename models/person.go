package models

type Person struct {
	NIC     string `json:"nic"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type IDChecked struct {
	Exists bool   `json:"exists"`
	NIC    string `json:"nic"`
}

type AddressChecked struct {
	NIC    string `json:"nic"`
	Exists bool   `json:"exists"`
}

type PoliceCheck struct {
	NIC   string `json:"nic"`
	Clear bool   `json:"clear"`
}

type StatusCheck struct {
	NIC   string `json:"nic"`
	Email string `json:"email"`
}

type AuthorizedUser struct {
	Name      string `json:"name"`
	Sub       string `json:"sub"`
	NIC       string `json:"nic"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Email     string `json:"username"`
}

type Check struct {
	Nic          string `json:"nic"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Idcheck      bool   `json:"idcheck"`
	Addresscheck bool   `json:"addresscheck"`
	Policecheck  bool   `json:"policecheck"`
}
