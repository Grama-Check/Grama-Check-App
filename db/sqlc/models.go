// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import ()

type Check struct {
	Nic          string `json:"nic"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Idcheck      bool   `json:"idcheck"`
	Addresscheck bool   `json:"addresscheck"`
	Policecheck  bool   `json:"policecheck"`
	Failed       bool   `json:"failed"`
}
