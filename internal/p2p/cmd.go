package p2p

import (
	"bufio"
	"fmt"
	"github.com/akbariandev/jumpy/internal/chain"
	"log"
	"os"
	"strings"
)

type Command string

const (
	LogCommand                Command = "log"
	TransactionCommand        Command = "transaction"
	CommitTransactionsCommand Command = "commit"
)

func (ps *PeerStream) readCli() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		inp, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		inp = strings.Replace(inp, "\n", "", -1)
		if strings.Contains(inp, ":") {
			scmd := strings.SplitN(inp, ":", 2)
			cmd := Command(scmd[0])
			data := []byte(scmd[1])
			cmd.run(ps, data)
		} else {
			cmd := Command(inp)
			cmd.run(ps, nil)
		}
	}
}

func (cmd Command) run(ps *PeerStream, data any) {
	switch cmd {
	case LogCommand:
		chain.PrintBlockChain()
	case TransactionCommand:
		addTransaction(ps, data)
	case CommitTransactionsCommand:
		commitTransaction(ps)
	default:
		fmt.Println("command not defined")
	}
}
