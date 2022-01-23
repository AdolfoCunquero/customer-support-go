package websocket

import (
	mdls "customer-support/models"
	mdb "customer-support/mongo"
	rdb "customer-support/redis"
	"log"
	"time"
)

func SaveMessage(msg mdls.Message) (mdls.Message, error) {
	msg.DateTime = time.Now()

	incidentId, errI := rdb.GetOrCreateIncident(msg.ClientId, "acunqueroc")

	if errI != nil {
		log.Println(errI.Error())
	}

	msg.IncidentId = incidentId
	msg.AgentId = "acunqueroc"

	result, err := mdb.SaveMessage(msg)
	if err != nil {
		log.Println(err.Error())
		return result, err
	}
	return result, nil
}
