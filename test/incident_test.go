package test

import (
	mdb "customer-support/mongo"
	rdb "customer-support/redis"
	"encoding/json"
	"fmt"
	"testing"
)

func TestSaveIncident(t *testing.T) {
	newInc, err := rdb.GetOrCreateIncident("cliente2", "adolfo")
	if err != nil {
		t.Error(err)

	}
	t.Log(newInc)
}

func TestCloseIncident(t *testing.T) {
	err := rdb.CloseIncident("cliente2")
	if err != nil {
		t.Error(err)
	}

	t.Log("finalizado")
}

func TestGetActiveIncidents(t *testing.T) {
	incidents, err := mdb.GetActiveIncidents("acunqueroc")
	if err != nil {
		t.Error(err)
	}

	b, err1 := json.Marshal(&incidents)
	if err1 != nil {
		t.Error(err)
	}

	fmt.Println(string(b))
	t.Log("FInalizado")
}

func TestGetActiveConversation(t *testing.T) {
	conversation, err := mdb.GetActiveConversation("61e12342d9388a555a40437a")
	if err != nil {
		t.Error(err)
	}

	b, err1 := json.Marshal(&conversation)
	if err1 != nil {
		t.Error(err)
	}

	fmt.Println(string(b))
	t.Log("FInalizado")
}
