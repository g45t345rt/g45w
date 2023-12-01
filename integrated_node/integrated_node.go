package integrated_node

import (
	"net"
	"time"

	"github.com/deroproject/derohe/blockchain"
	derodrpc "github.com/deroproject/derohe/cmd/derod/rpc"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/metrics"
	"github.com/deroproject/derohe/p2p"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
)

var Chain *blockchain.Blockchain
var RPCServer *derodrpc.RPCServer
var Running bool

func isTcpPortInUse(addr string) bool {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	defer conn.Close()

	return false
}

// func Start() error {
// 	if Running {
// 		return nil
// 	}

// 	rpcPort := config.Mainnet.RPC_Default_Port
// 	rpcAddr := fmt.Sprintf("127.0.0.1:%d", rpcPort)
// 	if isTcpPortInUse(rpcAddr) {
// 		return fmt.Errorf("rpc port (%d) already in use", rpcPort)
// 	}

// 	nodeDir := settings.IntegratedNodeDir

// 	runtime.MemProfileRate = 0

// 	globals.Arguments["--timeisinsync"] = false
// 	globals.Arguments["--p2p-bind"] = nil
// 	globals.Arguments["--add-exclusive-node"] = make([]string, 0)
// 	globals.Arguments["--node-tag"] = nil
// 	globals.Arguments["--clog-level"] = nil
// 	globals.Arguments["--sync-node"] = false
// 	globals.Arguments["--data-dir"] = nodeDir
// 	globals.Arguments["--add-priority-node"] = make([]string, 0)
// 	globals.Arguments["--rpc-bind"] = nil
// 	globals.Arguments["--integrator-address"] = nil
// 	globals.Arguments["--fastsync"] = true
// 	globals.Arguments["--socks-proxy"] = nil
// 	globals.Arguments["--min-peers"] = nil
// 	globals.Arguments["--max-peers"] = nil
// 	globals.Arguments["--getwork-bind"] = nil
// 	globals.Arguments["--prune-history"] = nil

// 	globals.InitializeLog(os.Stdout, io.Discard)
// 	globals.Initialize()

// 	params := make(map[string]interface{})

// 	chain, err := blockchain.Blockchain_Start(params)
// 	if err != nil {
// 		return err
// 	}

// 	Chain = chain
// 	params["chain"] = chain

// 	err = p2p.P2P_Init(params)
// 	if err != nil {
// 		return err
// 	}

// 	chain.P2P_Block_Relayer = func(cbl *block.Complete_Block, peerid uint64) {
// 		p2p.Broadcast_Block(cbl, peerid)
// 	}

// 	chain.P2P_MiniBlock_Relayer = func(mbl block.MiniBlock, peerid uint64) {
// 		p2p.Broadcast_MiniBlock(mbl, peerid)
// 	}

// 	RPCServer, err = derodrpc.RPCServer_Start(params)
// 	if err != nil {
// 		return err
// 	}

// 	globals.Cron.Start()
// 	Running = true
// 	return nil
// }

func Stop() {
	RPCServer.RPCServer_Stop()
	p2p.P2P_Shutdown() // does not close process_outgoing_connection? :(
	Chain.Shutdown()
	globals.Cron.Stop()
	metrics.Set.UnregisterAllMetrics()
	Running = false
}

type NodeStatus struct {
	Height          int64
	BestHeight      int64
	StableHeight    int64
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
				nodeStatus.Update()
				window.Invalidate()
				nodeStatus.isActive = false
			}
		}
	}()

	nodeStatus.Update()
	return nodeStatus
}

func (n *NodeStatus) Active() {
	n.isActive = true
}

func (n *NodeStatus) Update() {
	if Chain == nil {
		return
	}

	n.Height = Chain.Get_Height()
	bestHeight, _ := p2p.Best_Peer_Height()
	n.StableHeight = Chain.Get_Stable_Height()
	n.BestHeight = bestHeight
	//topo_height := chain.Load_TOPO_HEIGHT()

	n.MemCount = len(Chain.Mempool.Mempool_List_TX())
	n.RegCount = len(Chain.Regpool.Regpool_List_TX())

	//p2p.PeerList_Print()
	//n.PeerCount = p2p.Peer_Count()
	in, out := p2p.Peer_Direction_Count()
	n.PeerInCount = in
	n.PeerOutCount = out

	n.NetworkHashRate = Chain.Get_Network_HashRate()

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

	return nodedSize
}

func (n *NodeSize) Active() {
	n.isActive = true
}

func (n *NodeSize) update() {
	nodeDir := settings.IntegratedNodeDir
	size, _ := utils.GetFolderSize(nodeDir)
	n.Size = size
}
