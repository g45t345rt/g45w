package node_manager

import (
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/integrated_node"
	"github.com/g45t345rt/g45w/settings"
)

var CurrentNode *app_data.NodeConnection

func Load() error {
	endpoint := settings.App.NodeEndpoint
	if endpoint != "" {
		var nodeConn *app_data.NodeConnection
		integratedNodeConn := app_data.INTEGRATED_NODE_CONNECTION
		if endpoint == integratedNodeConn.Endpoint {
			nodeConn = &integratedNodeConn
		} else {
			nodeConn, err := app_data.GetNodeConnection(endpoint)
			if err != nil {
				return err
			}

			if nodeConn == nil {
				nodeConn = &app_data.NodeConnection{
					Name:     "",
					Endpoint: endpoint,
				}
			}
		}

		err := Connect(*nodeConn, false)
		if err != nil {
			return err
		}

		CurrentNode = nodeConn
	}

	return nil
}

func Connect(nodeConn app_data.NodeConnection, save bool) error {
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

	if integrated_node.Running &&
		settings.App.NodeEndpoint == app_data.INTEGRATED_NODE_CONNECTION.Endpoint {
		integrated_node.Stop()
	}

	CurrentNode = &nodeConn
	settings.App.NodeEndpoint = nodeConn.Endpoint

	if save {
		err := settings.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
