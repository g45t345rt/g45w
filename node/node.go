package node

import (
	"io"
	"os"
	"runtime"

	"github.com/deroproject/derohe/block"
	"github.com/deroproject/derohe/blockchain"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/settings"
)

type Node struct {
	Chain *blockchain.Blockchain
}

var Instance *Node

func NewNode() *Node {
	n := &Node{}
	Instance = n
	return n
}

func (n *Node) Start() error {
	nodeDir := settings.Instance.NodeDir

	runtime.MemProfileRate = 0
	globals.Arguments = make(map[string]interface{})

	globals.Arguments["--testnet"] = false
	globals.Arguments["--timeisinsync"] = false
	globals.Arguments["--p2p-bind"] = nil
	globals.Arguments["--add-exclusive-node"] = make([]string, 0)
	globals.Arguments["--node-tag"] = nil
	globals.Arguments["--clog-level"] = nil
	globals.Arguments["--debug"] = false
	globals.Arguments["--sync-node"] = false
	globals.Arguments["--data-dir"] = nodeDir
	globals.Arguments["--add-priority-node"] = make([]string, 0)
	globals.Arguments["--rpc-bind"] = nil
	globals.Arguments["--integrator-address"] = nil
	globals.Arguments["--flog-level"] = nil
	globals.Arguments["--fastsync"] = true
	globals.Arguments["--help"] = false
	globals.Arguments["--version"] = false
	globals.Arguments["--socks-proxy"] = nil
	globals.Arguments["--min-peers"] = nil
	globals.Arguments["--max-peers"] = nil
	globals.Arguments["--getwork-bind"] = nil
	globals.Arguments["--prune-history"] = nil
	globals.Arguments["--log-dir"] = nil

	globals.InitializeLog(os.Stdout, io.Discard)
	globals.Initialize()

	params := make(map[string]interface{})

	chain, err := blockchain.Blockchain_Start(params)
	if err != nil {
		return err
	}

	n.Chain = chain
	params["chain"] = chain

	err = p2p.P2P_Init(params)
	if err != nil {
		return err
	}

	chain.P2P_Block_Relayer = func(cbl *block.Complete_Block, peerid uint64) {
		p2p.Broadcast_Block(cbl, peerid)
	}

	chain.P2P_MiniBlock_Relayer = func(mbl block.MiniBlock, peerid uint64) {
		p2p.Broadcast_MiniBlock(mbl, peerid)
	}

	globals.Cron.Start()
	return nil
}
