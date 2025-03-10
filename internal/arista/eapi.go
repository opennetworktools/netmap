package arista

import (
	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
)

func GetNeighbors(hostname, username, password, enablePasswd string, port int) (*module.ShowLLDPNeighbors, error) {
	lldp := &module.ShowLLDPNeighbors{}

	node, err := goeapi.Connect("http", hostname, username, password, port)
	if err != nil {
		return lldp, err
	}
	node.EnableAuthentication(enablePasswd)

	handle, err := node.GetHandle("json")
	if err != nil {
		return lldp, err
	}
	handle.AddCommand(lldp)
	err = handle.Call()
	if err != nil {
		return lldp, err
	}

	return lldp, nil
}
