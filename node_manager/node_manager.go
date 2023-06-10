package node_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/g45t345rt/g45w/settings"
)

type NodeInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type NodeManager struct {
	UserNodes    map[string]NodeInfo
	TrustedNodes []NodeInfo
	Current      *NodeInfo // use integrated node if nil
}

var Instance *NodeManager

func NewNodeManager() *NodeManager {
	nodeManager := &NodeManager{
		UserNodes: make(map[string]NodeInfo),
		TrustedNodes: []NodeInfo{
			{ID: "deronfts", Name: "DeroNFTs", Host: "74.208.54.173", Port: 50404},
			{ID: "foundation", Name: "Foundation", Host: "ams.foundation.org", Port: 11011},
			{ID: "my_srv_cloud", Name: "MySrvCloud", Host: "213.171.208.37", Port: 18089},
			{ID: "derostats", Name: "DeroStats", Host: "163.172.26.245", Port: 10505},
		},
	}

	Instance = nodeManager
	return nodeManager
}

func (n *NodeManager) LoadNodes() error {
	nodeDir := settings.Instance.NodeDir

	path := filepath.Join(nodeDir, "list.json")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.WriteFile(path, []byte("{}"), os.ModeAppend)
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

	var nodes map[string]NodeInfo
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		return err
	}

	n.UserNodes = nodes
	return nil
}

func (n *NodeManager) AddNode(name string, host string, port uint) error {
	id := fmt.Sprint(time.Now().Unix())
	nodeInfo := NodeInfo{
		ID:   id,
		Name: name,
		Host: host,
		Port: port,
	}

	n.UserNodes[nodeInfo.ID] = nodeInfo
	return n.saveNodes()
}

func (n *NodeManager) DelNode(id string) error {
	delete(n.UserNodes, id)
	return n.saveNodes()
}

func (n *NodeManager) ClearNodes() error {
	n.UserNodes = make(map[string]NodeInfo)
	return n.saveNodes()
}

func (n *NodeManager) saveNodes() error {
	nodeDir := settings.Instance.NodeDir
	data, err := json.Marshal(n.UserNodes)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(nodeDir, "list.json"), data, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}
