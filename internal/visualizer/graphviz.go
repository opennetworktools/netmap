package visualizer

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/opennetworktools/netmap/internal/utils"

	"aqwari.net/xml/xmltree"
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

func SaveTopologyWithGraphviz(ctx context.Context, networkMap *NetworkMap, timestamp string, protocol string) error {
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
		node, err := graph.CreateNodeByName(device)
		if err != nil {
			fmt.Println("Error creating node:", err)
			continue
		}
		node.SetLabel(utils.TruncateString(device, 15))
		node.SetShape(cgraph.EllipseShape)
		node.SetStyle(cgraph.FilledNodeStyle)
		nodes[device] = node
	}

	// Create edges
	for device, edges := range networkMap.Devices {
		for _, edge := range edges {
			neighbor := edge.Neighbor
			if nodes[neighbor] == nil {
				node, err := graph.CreateNodeByName(neighbor)
				if err != nil {
					fmt.Println("Error creating node:", err)
					continue
				}
				node.SetLabel(utils.TruncateString(neighbor, 15))
				node.SetShape(cgraph.EllipseShape)
				node.SetStyle(cgraph.FilledNodeStyle)
				nodes[neighbor] = node
			}

			e, err := graph.CreateEdgeByName("", nodes[device], nodes[neighbor])
			if err != nil {
				fmt.Println("Error creating edge:", err)
				continue
			}
			if protocol == "ospf" {
				e.SetLabel(fmt.Sprintf("%s - %s", edge.LocalPort, edge.NeighborPort))
			} else {
				e.SetLabel(fmt.Sprintf("%s - %s", formatEdgeName(edge.LocalPort), formatEdgeName(edge.NeighborPort)))
			}
		}
	}

	// Dot file format
	// var buf bytes.Buffer
	// if err := g.Render(ctx, graph, "dot", &buf); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("\n%v", buf.String())

	// image.Image format
	outputImg, err := g.RenderImage(ctx, graph)
	if err != nil {
		return err
	}

	// encode to svg
	outputBuf := &bytes.Buffer{}
	err = g.Render(ctx, graph, graphviz.SVG, outputBuf)
	if err != nil {
		return err
	}

	doc, err := xmltree.Parse(outputBuf.Bytes())
	if err != nil {
		return err
	}

	outputSvg := xmltree.Marshal(doc)

	// MacOS Path: /Users/<USERNAME>/Library/Application Support/netmap
	appDir, err := utils.CreateDirectoryToSaveOutput()
	if err != nil {
		return err
	}

	// create filenames for svg and png
	baseFilename := "graph"
	svgFilename := fmt.Sprintf("%s_%s.svg", baseFilename, timestamp)
	svgFilenamePath := filepath.Join(appDir, svgFilename)
	pngFilename := fmt.Sprintf("%s_%s.png", baseFilename, timestamp)

	err = renderOutputToFilename(svgFilenamePath, outputSvg)
	if err != nil {
		return err
	}

	// encode to png
	outputImgBytes := make([]byte, 0)
	writer := bytes.NewBuffer(outputImgBytes)

	err = png.Encode(writer, outputImg)
	if err != nil {
		return err
	}

	err = renderOutputToFilename(filepath.Join(appDir, pngFilename), writer.Bytes())
	if err != nil {
		return err
	}

	fmt.Printf("\nExported graph to .json, .svg, and .png file formats: %s", appDir)
	fmt.Printf("\nTo view the SVG file, copy and paste the path in your favorite browser: %s", svgFilenamePath)

	// TODO: Create a ZIP file with .svg, .png, .jpg, dot files.

	// Save as PNG
	// filePath := "/Users/roopesh/Desktop/projects/network-mapper/network-topology.png"
	// if err := g.RenderFilename(ctx, graph, graphviz.PNG, filePath); err != nil {
	// 	fmt.Println("Error generating network topology PNG:", err)
	// 	return err
	// }
	// fmt.Println("Network topology image saved:", filePath)

	return nil
}

func renderOutputToFilename(filePath string, outData []byte) error {
	// creating a file
	// filePath := filepath.Join(appDir, "network-topology.png")
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// writing to a file
	_, err = outputFile.Write(outData)
	if err != nil {
		return err
	}

	return nil
}

func formatEdgeName(name string) string {
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		return parts[0]
	}
	return name
}
