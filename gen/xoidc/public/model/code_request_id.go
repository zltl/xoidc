//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type CodeRequestID struct {
	Code       string    `sql:"primary_key"`
	RequestID  uuid.UUID `sql:"primary_key"`
	CreateTime time.Time
}
