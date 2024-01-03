package node_manager

import (
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/settings"
)

var CurrentNode *app_db.NodeConnection

func Load() error {
	endpoint := settings.App.NodeEndpoint
	if endpoint != "" {
		var nodeConn *app_db.NodeConnection
		integratedNodeConn := app_db.INTEGRATED_NODE_CONNECTION
		if endpoint == integratedNodeConn.Endpoint {
			nodeConn = &integratedNodeConn
		} else {
			conn, err := app_db.GetNodeConnectionByEndpoint(endpoint)
			if err != nil {
				return err
			}

			nodeConn = &conn
		}

		err := Set(nodeConn, false)
		if err != nil {
			return err
		}

		CurrentNode = nodeConn
	}

	return nil
}

func Set(nodeConn *app_db.NodeConnection, save bool) error {
	integratedEndpoint := app_db.INTEGRATED_NODE_CONNECTION.Endpoint
	if nodeConn != nil {
		endpoint := nodeConn.Endpoint
		if nodeConn.Endpoint == integratedEndpoint {
			err := integrated_node.Start()
			if err != nil {
				return err
			}

			endpoint = "ws://127.0.0.1:10102/ws"
		}

		err := walletapi.Connect(endpoint)
		if err != nil {
			return err
		}

		settings.App.NodeEndpoint = nodeConn.Endpoint
	} else {
		go func() {
			rpcClient := walletapi.GetRPCClient()
			rpcClient.WS.Close()
			rpcClient.RPC.Close()
		}()

		settings.App.NodeEndpoint = ""
	}

	if integrated_node.Running &&
		settings.App.NodeEndpoint == integratedEndpoint {
		integrated_node.Stop()
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
