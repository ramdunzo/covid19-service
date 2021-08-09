package id_generator

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GetUniqId() string {
	randId := strings.Replace(uuid.Must(uuid.NewV4(), nil).String(), "-", "", -1)
	timeId := strings.Replace(uuid.Must(uuid.NewV1(), nil).String(), "-", "", -1)
	return timeId[:len(timeId)/2] + randId[len(randId)/2:]
}

func GetUUID() string {

	return uuid.NewV4().String()
}