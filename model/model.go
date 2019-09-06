package model

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

type DB struct {
	*mgo.Database
}

func genUUID() string {
	id := uuid.NewV4()
	return strings.Replace(fmt.Sprintf("%v", id), "-", "", -1)
}
