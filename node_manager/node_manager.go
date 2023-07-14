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

var CurrentNode string
var Nodes map[string]NodeConnection

const INTEGRATED_NODE_ID = "integrated"
const NODE_STATE_FILE = "node_state.json"

var TrustedNodes = map[string]NodeConnection{
	"deronfts":        {ID: "deronfts", Name: "DeroNFTs", Endpoint: "wss://node.deronfts.com/ws"},
	"my_srv_cloud":    {ID: "my_srv_cloud", Name: "MySrvCloud", Endpoint: "wss://dero-node.mysrv.cloud/ws"},
	"derostats":       {ID: "derostats", Name: "DeroStats", Endpoint: "ws://derostats.io:10102/ws"},
	"dero_foundation": {ID: "dero_foundation", Name: "Foundation", Endpoint: "ws://node.derofoundation.org:11012/ws"},
	"friendspool":     {ID: "friendspool", Name: "Friendspool", Endpoint: "ws://wallet.friendspool.club:10102/ws"},
}

func Load() error {
	nodeDir := settings.NodeDir

	err := os.MkdirAll(nodeDir, os.ModePerm)
	if err != nil {
		return err
	}

	path := filepath.Join(nodeDir, NODE_STATE_FILE)
	_, err = os.Stat(path)
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

	CurrentNode = nodeState.Current
	Nodes = nodeState.Nodes

	return nil
}

func AddNode(conn NodeConnection) error {
	id := fmt.Sprint(time.Now().Unix())
	conn.ID = id
	Nodes[id] = conn
	return saveState()
}

func EditNode(conn NodeConnection) error {
	Nodes[conn.ID] = conn
	return saveState()
}

func DelNode(id string) error {
	delete(Nodes, id)
	return saveState()
}

func ReloadTrustedNodes() error {
	for _, node := range TrustedNodes {
		Nodes[node.ID] = node
	}

	return saveState()
}

func ConnectNode(id string, save bool) error {
	if id != INTEGRATED_NODE_ID {
		_, ok := Nodes[id]
		if !ok {
			return fmt.Errorf("node id does not exists")
		}
	}

	if id == INTEGRATED_NODE_ID {
		err := integrated_node.Start()
		if err != nil {
			return err
		}

		localEndpoint := "ws://127.0.0.1:10102/ws"
		err = walletapi.Connect(localEndpoint)
		if err != nil {
			return err
		}
	} else {
		nodeConn := Nodes[id]
		err := walletapi.Connect(nodeConn.Endpoint)
		if err != nil {
			return err
		}

		if CurrentNode == INTEGRATED_NODE_ID {
			integrated_node.Stop()
		}
	}

	CurrentNode = id

	if save {
		err := saveState()
		if err != nil {
			return err
		}
	}

	return nil
}

func saveState() error {
	nodeDir := settings.NodeDir

	nodeState := NodeState{
		Nodes:   Nodes,
		Current: CurrentNode,
	}

	data, err := json.Marshal(nodeState)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(nodeDir, NODE_STATE_FILE), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
