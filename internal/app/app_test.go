package app

import (
	"testing"
	"time"
)

func TestApplication_Start(t *testing.T) {
	testData := []string{"data1", "data2", "data3", "data4"}
	app := &Application{}
	t.Log("App Started!")
	go app.Start(3, "")

	time.Sleep(2 * time.Second)
	if len(app.nodes) != 3 {
		t.Failed()
	}

	//commit some transaction
	node0 := app.nodes[0]
	node0.AddTransaction(testData[0])
	if err := node0.CommitTransaction(); err != nil {
		t.Fatal(err)
	}
	node0.AddTransaction(testData[1])
	if err := node0.CommitTransaction(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	node0Chain := node0.GetChain()
	if node0Chain.GetHeight() != 3 {
		t.Fatalf("expected %d got %d", 3, node0Chain.GetHeight())
	}
	node0B0Data := node0Chain.GetBlockAtIndex(1).Transaction[0].Data
	if node0B0Data != testData[0] {
		t.Fatalf("expected %s got %s", testData[0], node0B0Data)
	}
	if node0Chain.GetBlockAtIndex(1).Connections[0].PeerID != node0.Host.ID().String() {
		t.Fatalf("expected %s got %s", node0.Host.ID().String(), node0Chain.GetBlockAtIndex(1).Connections[0].PeerID)
	}
	if node0Chain.GetBlockAtIndex(1).Connections[1].PeerID == node0.Host.ID().String() {
		t.Fatalf("peer ID is equal to self ID")
	}
}
