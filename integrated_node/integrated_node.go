package integrated_node

import (
	"io"
	"os"
	"runtime"
	"time"

	"github.com/deroproject/derohe/block"
	"github.com/deroproject/derohe/blockchain"
	derodrpc "github.com/deroproject/derohe/cmd/derod/rpc"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
)

type IntegratedNode struct {
	Chain     *blockchain.Blockchain
	RPCServer *derodrpc.RPCServer
}

var Instance *IntegratedNode

func Instantiate() *IntegratedNode {
	Instance = &IntegratedNode{}
	return Instance
}

func (n *IntegratedNode) Start() error {
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

	n.RPCServer, err = derodrpc.RPCServer_Start(params)
	if err != nil {
		return err
	}

	globals.Cron.Start()
	return nil
}

func (n *IntegratedNode) Stop() {
	n.RPCServer.RPCServer_Stop()
	p2p.P2P_Shutdown()
	n.Chain.Shutdown()
	globals.Cron.Stop()
}

type NodeStatus struct {
	Height          int64
	BestHeight      int64
	MemCount        int
	RegCount        int
	PeerInCount     uint64
	PeerOutCount    uint64
	NetworkHashRate uint64

	TimeOffset    time.Duration
	TimeOffsetNTP time.Duration
	TimeOffsetP2P time.Duration

	isActive bool
}

func NewNodeStatus(d time.Duration) *NodeStatus {
	nodeStatus := &NodeStatus{isActive: false}
	ticker := time.NewTicker(d)

	window := app_instance.Window
	go func() {
		for range ticker.C {
			if nodeStatus.isActive {
				nodeStatus.update()
				window.Invalidate()
				nodeStatus.isActive = false
			}
		}
	}()

	nodeStatus.update()
	return nodeStatus
}

func (n *NodeStatus) Active() {
	n.isActive = true
}

func (n *NodeStatus) update() {
	if Instance == nil {
		return
	}

	chain := Instance.Chain
	if chain == nil {
		return
	}

	n.Height = chain.Get_Height()
	bestHeight, _ := p2p.Best_Peer_Height()
	n.BestHeight = bestHeight
	//topo_height := chain.Load_TOPO_HEIGHT()

	n.MemCount = len(chain.Mempool.Mempool_List_TX())
	n.RegCount = len(chain.Regpool.Regpool_List_TX())

	//p2p.PeerList_Print()
	//n.PeerCount = p2p.Peer_Count()
	in, out := p2p.Peer_Direction_Count()
	n.PeerInCount = in
	n.PeerOutCount = out

	n.NetworkHashRate = chain.Get_Network_HashRate()

	n.TimeOffset = globals.GetOffset().Round(time.Millisecond)
	n.TimeOffsetNTP = globals.GetOffsetNTP().Round(time.Millisecond)
	n.TimeOffsetP2P = globals.GetOffsetP2P().Round(time.Millisecond)
}

type NodeSize struct {
	Size int64

	isActive bool
}

func NewNodeSize(d time.Duration) *NodeSize {
	ticker := time.NewTicker(d)
	window := app_instance.Window
	nodedSize := &NodeSize{
		Size:     0,
		isActive: false,
	}

	go func() {
		for range ticker.C {
			if nodedSize.isActive {
				nodedSize.update()
				window.Invalidate()
				nodedSize.isActive = false
			}
		}
	}()

	nodedSize.update()
	return nodedSize
}

func (n *NodeSize) Active() {
	n.isActive = true
}

func (n *NodeSize) update() {
	nodeDir := settings.Instance.NodeDir
	size, _ := utils.GetFolderSize(nodeDir)
	n.Size = size
}
