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
	for _, d := range testData {
		node0.AddTransaction(d)
		if err := node0.CommitTransaction(); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(2 * time.Second)

	node0Chain := node0.GetChain()
	//check chain height
	if node0Chain.GetHeight() != 5 {
		t.Fatalf("expected %d got %d", 5, node0Chain.GetHeight())
	}

	//check blocks healthy
	node0B1Data := node0Chain.GetBlockAtIndex(1).Transaction[0].Data
	if node0B1Data != testData[0] {
		t.Fatalf("expected %s got %s", testData[0], node0B1Data)
	}
	node0B2Data := node0Chain.GetBlockAtIndex(2).Transaction[0].Data
	if node0B2Data != testData[1] {
		t.Fatalf("expected %s got %s", testData[1], node0B2Data)
	}
	node0B3Data := node0Chain.GetBlockAtIndex(3).Transaction[0].Data
	if node0B3Data != testData[2] {
		t.Fatalf("expected %s got %s", testData[2], node0B3Data)
	}
	node0B4Data := node0Chain.GetBlockAtIndex(4).Transaction[0].Data
	if node0B4Data != testData[3] {
		t.Fatalf("expected %s got %s", testData[3], node0B4Data)
	}

	//check block connections
	if node0Chain.GetBlockAtIndex(1).Connections[0].PeerID != node0.Host.ID().String() {
		t.Fatalf("expected %s got %s", node0.Host.ID().String(), node0Chain.GetBlockAtIndex(1).Connections[0].PeerID)
	}
	if node0Chain.GetBlockAtIndex(1).Connections[1].PeerID == node0.Host.ID().String() {
		t.Fatalf("peer ID is equal to self ID")
	}
}
