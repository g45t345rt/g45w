package node_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/settings"
)

type NodeConnection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type NodeState struct {
	Nodes   map[string]NodeConnection `json:"nodes"`
	Current string                    `json:"current"`
}

const INTEGRATED_NODE_ID = "integrated"
const NODE_STATE_FILE = "node_state.json"

type NodeManager struct {
	NodeState NodeState
}

var Instance *NodeManager

var TrustedNodes = map[string]NodeConnection{
	"deronfts":     {ID: "deronfts", Name: "DeroNFTs", Endpoint: "wss://node.deronfts.com/ws"},
	"my_srv_cloud": {ID: "my_srv_cloud", Name: "MySrvCloud", Endpoint: "wss://dero-node.mysrv.cloud/ws"},
	"derostats":    {ID: "derostats", Name: "DeroStats", Endpoint: "ws://derostats.io:10102/ws"},
}

func Instantiate() *NodeManager {
	Instance = &NodeManager{}
	return Instance
}

func (n *NodeManager) Load() error {
	nodeDir := settings.Instance.NodeDir

	path := filepath.Join(nodeDir, NODE_STATE_FILE)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		nodeState := NodeState{
			Nodes: TrustedNodes,
		}

		data, err := json.Marshal(nodeState)
		if err != nil {
			return err
		}

		err = os.WriteFile(path, data, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var nodeState NodeState
	err = json.Unmarshal(data, &nodeState)
	if err != nil {
		return err
	}

	n.NodeState = nodeState
	n.SelectNode(n.NodeState.Current, false)
	return nil
}

func (n *NodeManager) AddNode(conn NodeConnection) error {
	id := fmt.Sprint(time.Now().Unix())
	conn.ID = id
	n.NodeState.Nodes[id] = conn
	return n.saveState()
}

func (n *NodeManager) EditNode(conn NodeConnection) error {
	n.NodeState.Nodes[conn.ID] = conn
	return n.saveState()
}

func (n *NodeManager) DelNode(id string) error {
	delete(n.NodeState.Nodes, id)
	return n.saveState()
}

func (n *NodeManager) ReloadTrustedNodes() error {
	for _, node := range TrustedNodes {
		n.NodeState.Nodes[node.ID] = node
	}

	return n.saveState()
}

func (n *NodeManager) SelectNode(id string, save bool) error {
	if id != INTEGRATED_NODE_ID {
		_, ok := n.NodeState.Nodes[id]
		if !ok {
			return fmt.Errorf("node id does not exists")
		}
	}

	if id == INTEGRATED_NODE_ID {
		err := integrated_node.Instance.Start()
		if err != nil {
			return err
		}

		err = walletapi.Connect("ws://127.0.0.1:10102/ws")
		if err != nil {
			return err
		}
	} else {
		nodeConn := n.NodeState.Nodes[id]
		err := walletapi.Connect(nodeConn.Endpoint)
		if err != nil {
			return err
		}

		if n.NodeState.Current == INTEGRATED_NODE_ID {
			integrated_node.Instance.Stop()
		}
	}

	n.NodeState.Current = id

	if save {
		err := n.saveState()
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NodeManager) saveState() error {
	nodeDir := settings.Instance.NodeDir
	data, err := json.Marshal(n.NodeState)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(nodeDir, NODE_STATE_FILE), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
