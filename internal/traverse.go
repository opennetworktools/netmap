package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opennetworktools/netmap/internal/utils"
	"github.com/opennetworktools/netmap/internal/visualizer"

	"github.com/opennetworktools/netmap/internal/arista"
)

func Traverse(hostname, username, password string) {
	ctx := context.Background()

	hostname = normalizeHostname(hostname)

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
			neighborDevice := normalizeHostname(obj.NeighborDevice)
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

	// Generate timestamp for filename (.json, .svg, .png)
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Exporting networkMap to a JSON file
	err := utils.SaveStructAsJson("graph", networkMap, timestamp)
	if err != nil {
		fmt.Println("Error writing the file: ", err)
		return
	}

	// Creating a graph and exporting it to a PNG file using go-graphviz
	err = visualizer.SaveTopologyWithGraphviz(ctx, networkMap, timestamp, "lldp")
	if err != nil {
		fmt.Println("Error rendering the topology: ", err)
		return
	}
}

// Sometimes the neighbor or peer device sees the FQDN instead of hostname, which is false.
// For eg. Say the hostname is "ok270" and the neigbor device sees this as "ok270.aristanetworks.com"
// then an extra API call will be triggered as the hostname and FQDN didn't match.
func normalizeHostname(hostname string) string {
	parts := strings.Split(hostname, ".")
	return parts[0]
}
