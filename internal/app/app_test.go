package app

import (
	"testing"
	"time"
)

func TestApplication_SingleCommitOnSingleNode(t *testing.T) {
	testData := "data1"
	app := &Application{}
	t.Log("App Started!")
	go app.Start(3, "")

	time.Sleep(2 * time.Second)
	if len(app.nodes) != 3 {
		t.Failed()
	}

	//commit some transaction
	node0 := app.nodes[0]
	node0.AddTransaction(testData)
	if err := node0.CommitTransaction(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	node0Chain := node0.GetChain()
	//check chain height
	if node0Chain.GetHeight() != 2 {
		t.Fatalf("expected %d got %d", 2, node0Chain.GetHeight())
	}

	//check blocks healthy
	node0B1Data := node0Chain.GetBlockAtIndex(1).Transaction[0].Data
	if node0B1Data != testData {
		t.Fatalf("expected %s got %s", testData, node0B1Data)
	}

	//check block connections
	if node0Chain.GetBlockAtIndex(1).Connections[0].PeerID != node0.Host.ID().String() {
		t.Fatalf("expected %s got %s", node0.Host.ID().String(), node0Chain.GetBlockAtIndex(1).Connections[0].PeerID)
	}
	if node0Chain.GetBlockAtIndex(1).Connections[1].PeerID == node0.Host.ID().String() {
		t.Fatalf("peer ID is equal to self ID")
	}
}

func TestApplication_MultiCommitOnSingleNode(t *testing.T) {
	testData := []string{"data1", "data2", "data3", "data4"}
	app := &Application{}
	t.Log("App Started!")
	go app.Start(3, "")

	time.Sleep(2 * time.Second)
	if len(app.nodes) != 3 {
		t.Failed()
	}

	//commit some transaction
	n := app.nodes[0]
	for _, d := range testData {
		n.AddTransaction(d)
		if err := n.CommitTransaction(); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(5 * time.Second)

	c := n.GetChain()
	//check chain height
	if c.GetHeight() != 5 {
		t.Fatalf("expected %d got %d", 5, c.GetHeight())
	}

	//check blocks healthy
	for i, tt := range testData {
		b := c.GetBlockAtIndex(i + 1) // 0 == genesis
		d := b.Transaction[0].Data
		if d != tt {
			t.Fatalf("expected %s got %s", tt, d)
		}
		if b.Connections[0].PeerID != n.Host.ID().String() {
			t.Fatalf("expected %s got %s", n.Host.ID().String(), b.Connections[0].PeerID)
		}
		if b.Connections[1].PeerID == n.Host.ID().String() {
			t.Fatalf("peer ID is equal to self ID")
		}
	}
}

func TestApplication_MultiCommitOnMultiNode(t *testing.T) {
	testData := []string{"data1", "data2", "data3", "data4"}
	app := &Application{}
	t.Log("App Started!")
	numOfNodes := 5
	go app.Start(numOfNodes, "")

	time.Sleep(2 * time.Second)
	if len(app.nodes) != numOfNodes {
		t.Failed()
	}

	//commit some transaction
	t.Parallel()
	for numOfNodes > 0 {
		n := app.nodes[numOfNodes-1]
		for _, d := range testData {
			n.AddTransaction(d)
			if err := n.CommitTransaction(); err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(10 * time.Second)

		c := n.GetChain()
		//check chain height
		if c.GetHeight() != 5 {
			t.Fatalf("expected %d got %d", 5, c.GetHeight())
		}

		//check blocks healthy
		for i, tt := range testData {
			b := c.GetBlockAtIndex(i + 1) // 0 == genesis
			d := b.Transaction[0].Data
			if d != tt {
				t.Fatalf("expected %s got %s", tt, d)
			}
			if b.Connections[0].PeerID != n.Host.ID().String() {
				t.Fatalf("expected %s got %s", n.Host.ID().String(), b.Connections[0].PeerID)
			}
			if b.Connections[1].PeerID == n.Host.ID().String() {
				t.Fatalf("peer ID is equal to self ID")
			}
		}

		numOfNodes--
		go n.GetChain().PrintBlockChain()
	}

}
