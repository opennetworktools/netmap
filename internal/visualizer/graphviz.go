package visualizer

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"opennetworktools/netmap/internal/utils"
	"os"
	"path/filepath"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type NetworkMap struct {
	Devices map[string][]Edge
}

type Edge struct {
	LocalPort    string
	Neighbor     string
	NeighborPort string
}

func SaveTopologyWithGraphviz(ctx context.Context, networkMap *NetworkMap) error {
	g, err := graphviz.New(ctx)
	if err != nil {
		return err
	}
	graph, err := g.Graph(graphviz.WithDirectedType(cgraph.UnDirected))
	if err != nil {
		return err
	}
	defer g.Close()

	// Left to Right Direction
	graph.SetRankDir(graphviz.LRRank)

	nodes := make(map[string]*cgraph.Node)

	// Create nodes
	for device := range networkMap.Devices {
		node, err := graph.CreateNodeByName(utils.TruncateString(device, 8))
		if err != nil {
			fmt.Println("Error creating node:", err)
			continue
		}
		node.SetShape(cgraph.EllipseShape)
		node.SetStyle(cgraph.FilledNodeStyle)
		nodes[device] = node
	}

	// Create edges
	for device, edges := range networkMap.Devices {
		for _, edge := range edges {
			neighbor := edge.Neighbor
			if nodes[neighbor] == nil {
				node, err := graph.CreateNodeByName(utils.TruncateString(neighbor, 8))
				if err != nil {
					fmt.Println("Error creating node:", err)
					continue
				}
				node.SetShape(cgraph.EllipseShape)
				node.SetStyle(cgraph.FilledNodeStyle)
				nodes[neighbor] = node
			}

			e, err := graph.CreateEdgeByName("", nodes[device], nodes[neighbor])
			if err != nil {
				fmt.Println("Error creating edge:", err)
				continue
			}
			e.SetLabel(fmt.Sprintf("%s ↔ %s", edge.LocalPort, edge.NeighborPort))
		}
	}

	// Dot file format
	// var buf bytes.Buffer
	// if err := g.Render(ctx, graph, "dot", &buf); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(buf.String())

	// image.Image format
	outputImg, err := g.RenderImage(ctx, graph)
	if err != nil {
		return err
	}

	// encode
	out := make([]byte, 0)
	writer := bytes.NewBuffer(out)

	err = png.Encode(writer, outputImg)
	if err != nil {
		return err
	}

	appDir, err := utils.CreateDirectoryToSaveOutput()
	if err != nil {
		return err
	}

	// creating a file
	filePath := filepath.Join(appDir, "network-topology.png")
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// writing to a file
	_, err = outputFile.Write(writer.Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("\nOutput path: %s", appDir)

	// Save as PNG
	// filePath := "/Users/roopesh/Desktop/projects/network-mapper/network-topology.png"
	// if err := g.RenderFilename(ctx, graph, graphviz.PNG, filePath); err != nil {
	// 	fmt.Println("Error generating network topology PNG:", err)
	// 	return err
	// }
	// fmt.Println("Network topology image saved:", filePath)

	return nil
}
