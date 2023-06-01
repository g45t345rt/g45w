package node

import (
	"io"
	"os"
	"runtime"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/deroproject/derohe/block"
	"github.com/deroproject/derohe/blockchain"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
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

type NodeStatus struct {
	Height          int64
	BestHeight      int64
	MemCount        int
	RegCount        int
	PeerCount       uint64
	NetworkHashRate uint64

	TimeOffset    time.Duration
	TimeOffsetNTP time.Duration
	TimeOffsetP2P time.Duration

	startTime      time.Time
	updateInterval time.Duration
}

func NewNodeStatus(updateInterval time.Duration) *NodeStatus {
	return &NodeStatus{
		startTime:      time.Now().Add(-updateInterval),
		updateInterval: updateInterval,
	}
}

func (n *NodeStatus) Update(gtx layout.Context) {
	elapsed := gtx.Now.Sub(n.startTime)

	if elapsed < n.updateInterval {
		return
	}

	n.startTime = gtx.Now
	op.InvalidateOp{}.Add(gtx.Ops)

	chain := Instance.Chain
	n.Height = chain.Get_Height()
	bestHeight, _ := p2p.Best_Peer_Height()
	n.BestHeight = bestHeight
	//topo_height := chain.Load_TOPO_HEIGHT()

	n.MemCount = len(chain.Mempool.Mempool_List_TX())
	n.RegCount = len(chain.Regpool.Regpool_List_TX())

	//p2p.PeerList_Print()
	n.PeerCount = p2p.Peer_Count()
	//inc, out := p2p.Peer_Direction_Count()
	n.NetworkHashRate = chain.Get_Network_HashRate()

	n.TimeOffset = globals.GetOffset().Round(time.Millisecond)
	n.TimeOffsetNTP = globals.GetOffsetNTP().Round(time.Millisecond)
	n.TimeOffsetP2P = globals.GetOffsetP2P().Round(time.Millisecond)
}

type NodeSize struct {
	Size int64

	startTime      time.Time
	updateInterval time.Duration
}

func NewNodeSize(updateInterval time.Duration) *NodeSize {
	return &NodeSize{
		Size:           0,
		startTime:      time.Now().Add(-updateInterval),
		updateInterval: updateInterval,
	}
}

func (n *NodeSize) Update(gtx layout.Context) {
	elapsed := gtx.Now.Sub(n.startTime)

	if elapsed < n.updateInterval {
		return
	}

	n.startTime = gtx.Now
	op.InvalidateOp{}.Add(gtx.Ops)

	nodeDir := settings.Instance.NodeDir
	size, _ := utils.GetFolderSize(nodeDir)
	n.Size = size
}
