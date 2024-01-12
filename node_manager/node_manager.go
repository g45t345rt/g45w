package node_manager

import (
	"strconv"

	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/settings"
)

var CurrentNode *app_db.NodeConnection

func Load() error {
	nodeIdString := settings.App.NodeSelect
	if nodeIdString != "" {
		id, err := strconv.ParseInt(nodeIdString, 10, 64)
		if err != nil {
			return err
		}

		var nodeConn *app_db.NodeConnection
		integratedNodeConn := app_db.INTEGRATED_NODE_CONNECTION
		localNodeConn := app_db.LOCAL_NODE_CONNECTION
		switch id {
		case integratedNodeConn.ID:
			nodeConn = &integratedNodeConn
		case localNodeConn.ID:
			nodeConn = &localNodeConn
		default:
			conn, err := app_db.GetNodeConnection(id)
			if err != nil {
				return err
			}

			nodeConn = &conn
		}

		err = Set(nodeConn, false)
		if err != nil {
			return err
		}

		CurrentNode = nodeConn
	}

	return nil
}

func Set(nodeConn *app_db.NodeConnection, save bool) error {
	if nodeConn != nil {
		if nodeConn.Integrated {
			err := integrated_node.Start()
			if err != nil {
				return err
			}
		}

		err := walletapi.Connect(nodeConn.Endpoint)
		if err != nil {
			return err
		}

		settings.App.NodeSelect = strconv.FormatInt(nodeConn.ID, 10)
		if integrated_node.Running && !nodeConn.Integrated {
			integrated_node.Stop()
		}
	} else {
		go func() {
			rpcClient := walletapi.GetRPCClient()
			rpcClient.WS.Close()
			rpcClient.RPC.Close()
		}()

		settings.App.NodeSelect = ""
		if integrated_node.Running {
			integrated_node.Stop()
		}
	}

	CurrentNode = nodeConn

	if save {
		err := settings.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
