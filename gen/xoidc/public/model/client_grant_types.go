//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type ClientGrantTypes struct {
	ClientID  int64  `sql:"primary_key"`
	GrantType string `sql:"primary_key"`
}