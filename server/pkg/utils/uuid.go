package utils

import (
	"log"

	"github.com/google/uuid"
)

func GetUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		log.Panicln("获取uuid失败", err.Error())
	}
	return id.String()
}
