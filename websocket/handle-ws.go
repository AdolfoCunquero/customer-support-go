package websocket

import (
	mdls "customer-support/models"
	mng "customer-support/mongo"
	"log"
	"time"
)

func SaveMessage(msg mdls.Message) error {
	msg.DateTime = time.Now()
	err := mng.SaveMessage(msg)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
