package internal

import (
	"context"
	"fmt"
	"opennetworktools/netmap/internal/arista"
	"opennetworktools/netmap/internal/utils"
	"opennetworktools/netmap/internal/visualizer"
)

func Traverse(hostname, username, password string) {
	ctx := context.Background()

	// hostname := "roi460"
	// username := "admin"
	// password := ""
	enablePasswd := ""
	port := 80

	networkMap := &visualizer.NetworkMap{
		Devices: make(map[string][]visualizer.Edge),
	}

	// Queue for traversal
	pendingDevices := []string{hostname}
	visited := make(map[string]bool)

	for len(pendingDevices) > 0 {
		device := pendingDevices[0]
		pendingDevices = pendingDevices[1:]

		// Skip if already visited
		if visited[device] {
			continue
		}
		visited[device] = true

		fmt.Printf("\nDiscovering neighbors for %s...", device)

		// Get LLDP neighbors
		neighbors, err := arista.GetNeighbors(device, username, password, enablePasswd, port)
		if err != nil {
			fmt.Printf("\n%s", err.Error())
			continue
		}

		// Store device in topology (ensure it's initialized)
		if _, exists := networkMap.Devices[device]; !exists {
			networkMap.Devices[device] = []visualizer.Edge{}
		}

		fmt.Printf(" Found %d for %s!", len(neighbors.LLDPNeighbors), device)

		// Process each neighbor
		for _, obj := range neighbors.LLDPNeighbors {
			neighborDevice := obj.NeighborDevice
			// fmt.Printf("\n - Found neighbor: %s", neighborDevice)

			// Store the connection in the network map
			networkMap.Devices[device] = append(networkMap.Devices[device], visualizer.Edge{
				LocalPort:    obj.Port,
				Neighbor:     neighborDevice,
				NeighborPort: obj.NeighborPort,
			})

			// Add neighbor to queue if not visited
			if !visited[neighborDevice] {
				pendingDevices = append(pendingDevices, neighborDevice)
			}
		}
	}

	// Exporting networkMap to a JSON file
	err := utils.SaveStructAsJson("networkMap", networkMap)
	if err != nil {
		fmt.Println("Error writing the file: ", err)
		return
	}

	// Creating a graph and exporting it to a PNG file using go-graphviz
	err = visualizer.SaveTopologyWithGraphviz(ctx, networkMap)
	if err != nil {
		fmt.Println("Error rendering the topology: ", err)
		return
	}
}
