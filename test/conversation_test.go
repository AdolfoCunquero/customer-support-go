package test

import (
	rdb "customer-support/redis"
	"testing"
)

func TestSaveConversation(t *testing.T) {
	newConv, err := rdb.GetOrCreateConversation("cliente2")
	if err != nil {
		t.Error(err)
	}
	t.Log(newConv)
}

func TestCloseConversation(t *testing.T) {
	err := rdb.CloseConversation("cliente2")
	if err != nil {
		t.Error(err)
	}
	t.Log("Finalizado")

}
