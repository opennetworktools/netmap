package internal

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aristanetworks/goeapi"
	"github.com/opennetworktools/netmap/internal/utils"
	"github.com/opennetworktools/netmap/internal/visualizer"
)

type ShowIPOSPF struct {
	Vrfs map[string]VRF `json:"vrfs"`
}

type VRF struct {
	InstList map[string]OSPFInstance `json:"instList"`
}

type OSPFInstance struct {
	Asbr                         bool                   `json:"asbr"`
	LsaInformation               LsaInformation         `json:"lsaInformation"`
	GracefulRestartInfo          GracefulRestartInfo    `json:"gracefulRestartInfo"`
	ExternalLsaInformation       ExternalLsaInformation `json:"externalLsaInformation"`
	InstanceId                   int                    `json:"instanceId"`
	SpfInformation               SpfInformation         `json:"spfInformation"`
	ReferenceBandwidth           int                    `json:"referenceBandwidth"`
	OpaqueLsaInformation         OpaqueLsaInformation   `json:"opaqueLsaInformation"`
	TunnelRoutesEnabled          bool                   `json:"tunnelRoutesEnabled"`
	NumNormalAreas               int                    `json:"numNormalAreas"`
	MaxLsaInformation            MaxLsaInformation      `json:"maxLsaInformation"`
	EcmpMaximumNexthops          int                    `json:"ecmpMaximumNexthops"`
	NumBackboneNeighbors         int                    `json:"numBackboneNeighbors"`
	AdjacencyExchangeStartThresh int                    `json:"adjacencyExchangeStartThreshold"`
	AreaList                     map[string]Area        `json:"areaList"`
	NumAreas                     int                    `json:"numAreas"`
	Abr                          bool                   `json:"abr"`
	RouterId                     string                 `json:"routerId"`
	NumStubAreas                 int                    `json:"numStubAreas"`
	ShutDown                     bool                   `json:"shutDown"`
	NumNssaAreas                 int                    `json:"numNssaAreas"`
	LsaRetransmissionThreshold   int                    `json:"lsaRetransmissionThreshold"`
	FloodPacing                  int                    `json:"floodPacing"`
	Vrf                          string                 `json:"vrf"`
}

type LsaInformation struct {
	NumLsa             int `json:"numLsa"`
	LsaStartInterval   int `json:"lsaStartInterval"`
	LsaMaxWaitInterval int `json:"lsaMaxWaitInterval"`
	LsaHoldInterval    int `json:"lsaHoldInterval"`
	LsaArrivalInterval int `json:"lsaArrivalInterval"`
}

type GracefulRestartInfo struct {
	GracefulRestart       bool    `json:"gracefulRestart"`
	PlannedOnly           bool    `json:"plannedOnly"`
	HelperLooseLsaCheck   bool    `json:"helperLooseLsaCheck"`
	GracePeriod           float64 `json:"gracePeriod"`
	LastExitReason        string  `json:"lastExitReason"`
	LastExitTime          *string `json:"lastExitTime"` // nullable
	State                 string  `json:"state"`
	RestartExpirationTime *string `json:"restartExpirationTime"` // nullable
	LastRestartReason     string  `json:"lastRestartReason"`
	RestartDuration       *string `json:"restartDuration"` // nullable
	HelperMode            bool    `json:"helperMode"`
}

type ExternalLsaInformation struct {
	AsExternalCksum int `json:"asExternalCksum"`
	TreeSize        int `json:"treeSize"`
}

type SpfInformation struct {
	SpfStartInterval    int `json:"spfStartInterval"`
	SpfMaxWaitInterval  int `json:"spfMaxWaitInterval"`
	SpfCurrHoldInterval int `json:"spfCurrHoldInterval"`
	SpfHoldInterval     int `json:"spfHoldInterval"`
	SpfInterval         int `json:"spfInterval"`
	LastSpf             int `json:"lastSpf"`
	NextSpf             int `json:"nextSpf"`
}

type OpaqueLsaInformation struct {
	OpaqueTreeSize int  `json:"opaqueTreeSize"`
	OpaqueCksum    int  `json:"opaqueCksum"`
	Opaque         bool `json:"opaque"`
}

type MaxLsaInformation struct {
	MaxLsaCurrentIgnoreCount int  `json:"maxLsaCurrentIgnoreCount"`
	MaxLsaThreshold          int  `json:"maxLsaThreshold"`
	MaxLsaIgnoring           bool `json:"maxLsaIgnoring"`
	MaxLsaAllowedIgnoreCount int  `json:"maxLsaAllowedIgnoreCount"`
	MaxLsaWarningOnly        bool `json:"maxLsaWarningOnly"`
	MaxLsa                   int  `json:"maxLsa"`
	MaxLsaResetTime          int  `json:"maxLsaResetTime"`
	MaxLsaIgnoreTime         int  `json:"maxLsaIgnoreTime"`
}

type Area struct {
	NormalArea               bool                     `json:"normalArea"`
	LsaInformation           ExternalLsaInformation   `json:"lsaInformation"`
	AreaId                   string                   `json:"areaId"`
	TeEnabled                bool                     `json:"teEnabled"`
	SpfCount                 int                      `json:"spfCount"`
	OpaqueLsaInformation     OpaqueLsaInformation     `json:"opaqueLsaInformation"`
	NumIntf                  int                      `json:"numIntf"`
	RangeList                map[string]interface{}   `json:"rangeList"` // assuming empty map
	StubArea                 bool                     `json:"stubArea"`
	OpaqueAreaLsaInformation OpaqueAreaLsaInformation `json:"opaqueAreaLsaInformation"`
	AreaFiltersConfigured    bool                     `json:"areaFiltersConfigured"`
}

type OpaqueAreaLsaInformation struct {
	OpaqueTreeSize int `json:"opaqueTreeSize"`
	OpaqueCksumSum int `json:"opaqueCksumSum"`
}

func (s *ShowIPOSPF) GetCmd() string {
	return "show ip ospf"
}

// Neighbor types

type ShowIPOSPFNeighborDetail struct {
	Vrfs map[string]VRFNeighborInfo
}

type VRFNeighborInfo struct {
	InstList map[string]OSPFNeighborInstance `json:"instList"`
}

type OSPFNeighborInstance struct {
	OSPFNeighborEntries []OSPFNeighborEntry `json:"ospfNeighborEntries"`
}

type OSPFNeighborEntry struct {
	RouterId         string             `json:"routerId"`
	Priority         int                `json:"priority"`
	DrState          string             `json:"drState"`
	InterfaceName    string             `json:"interfaceName"`
	AdjacencyState   string             `json:"adjacencyState"`
	Inactivity       float64            `json:"inactivity"`
	InterfaceAddress string             `json:"interfaceAddress"`
	Options          OSPFOptions        `json:"options"`
	Details          OSPFNeighborDetail `json:"details"`
}

type OSPFOptions struct {
	MultitopologyCapability   bool `json:"multitopologyCapability"`
	DoNotUseInRouteCalc       bool `json:"doNotUseInRouteCalc"`
	DemandCircuitsSupport     bool `json:"demandCircuitsSupport"`
	NssaCapability            bool `json:"nssaCapability"`
	ExternalRoutingCapability bool `json:"externalRoutingCapability"`
	OpaqueLsaSupport          bool `json:"opaqueLsaSupport"`
	LinkLocalSignaling        bool `json:"linkLocalSignaling"`
	MulticastCapability       bool `json:"multicastCapability"`
}

type OSPFNeighborDetail struct {
	AreaId                 string  `json:"areaId"`
	BfdState               string  `json:"bfdState"`
	BackupDesignatedRouter string  `json:"backupDesignatedRouter"`
	StateTime              float64 `json:"stateTime"`
	RetxCount              int     `json:"retransmissionCount"`
	GrLastRestartTime      *string `json:"grLastRestartTime"` // nullable
	DesignatedRouter       string  `json:"designatedRouter"`
	GrHelperTimer          *string `json:"grHelperTimer"` // nullable
	NumberOfStateChanges   int     `json:"numberOfStateChanges"`
	BfdRequestSent         bool    `json:"bfdRequestSent"`
	GrNumAttempts          int     `json:"grNumAttempts"`
	InactivityDefers       int     `json:"inactivityDefers"`
}

func (s *ShowIPOSPFNeighborDetail) GetCmd() string {
	return "show ip ospf neighbor detail"
}

// Database Detail

type ShowIPOSPFDatabaseDetail struct {
	VRFs map[string]Vrf `json:"vrfs"`
}

type Vrf struct {
	InstList map[string]Instance `json:"instList"`
}

type Instance struct {
	Areas map[string]OSPFArea `json:"areas"`
}

type OSPFArea struct {
	AreaDatabase []AreaDatabaseEntry `json:"areaDatabase"`
}

type AreaDatabaseEntry struct {
	AreaLsas []AreaLsa `json:"areaLsas"`
}

type AreaLsa struct {
	LsaType           string          `json:"lsaType"`
	AdvertisingRouter string          `json:"advertisingRouter"`
	LinkStateID       string          `json:"linkStateId"`
	OspfRouterLsa     *OspfRouterLsa  `json:"ospfRouterLsa,omitempty"`
	OspfNetworkLsa    *OspfNetworkLsa `json:"ospfNetworkLsa,omitempty"`
}

type OspfRouterLsa struct {
	NumRtrLinks    int             `json:"numRtrLinks"`
	RouterLsaLinks []RouterLsaLink `json:"routerLsaLinks"`
}

type RouterLsaLink struct {
	LinkType string `json:"linkType"`
	LinkData string `json:"linkData"`
	LinkID   string `json:"linkId"`
	Metric   int    `json:"metric"`
	NumTos   int    `json:"numTos"`
}

type OspfNetworkLsa struct {
	NetworkMask     string   `json:"networkMask"`
	AttachedRouters []string `json:"attachedRouters"`
}

func (s *ShowIPOSPFDatabaseDetail) GetCmd() string {
	return "show ip ospf database detail"
}

// For graphing
type OSPFRouter struct {
	RouterID string
	Stubs    []Stub
	Transits []string
	Links    []Link
}

type Stub struct {
	Subnet string
	Mask   string
	Metric int
}

type Link struct {
	Neighbor  string
	IPAddress string
	Metric    int
}

func Ospf(hostname, username, password string) {
	ctx := context.Background()

	fmt.Println("OSPF")
	port := 80
	node, err := goeapi.Connect("http", hostname, username, password, port)
	if err != nil {
		fmt.Println(err.Error())
	}

	// showIPOSPF := &ShowIPOSPF{}
	// handle, _ := node.GetHandle("json")
	// handle.AddCommand(showIPOSPF)
	// if err := handle.Call(); err != nil {
	// 	panic(err)
	// }

	// showIPOSPFNeighborDetail := &ShowIPOSPFNeighborDetail{}
	// handle, _ = node.GetHandle("json")
	// handle.AddCommand(showIPOSPFNeighborDetail)
	// if err := handle.Call(); err != nil {
	// 	panic(err)
	// }

	showIPOSPFDatabaseDetail := &ShowIPOSPFDatabaseDetail{}
	handle, _ := node.GetHandle("json")
	handle.AddCommand(showIPOSPFDatabaseDetail)
	if err := handle.Call(); err != nil {
		panic(err)
	}

	// fmt.Println(showIPOSPFDatabaseDetail)

	// OSPFRouters := []OSPFRouter{}

	type OSPFRoutersMap struct {
		Routers map[string]OSPFRouter
	}

	ospfRoutersMap := &OSPFRoutersMap{
		Routers: make(map[string]OSPFRouter),
	}

	for _, inst := range showIPOSPFDatabaseDetail.VRFs {
		for _, area := range inst.InstList {
			for _, areaDb := range area.Areas {
				for _, lsa := range areaDb.AreaDatabase {
					// fmt.Println(lsa)
					for _, areaLsa := range lsa.AreaLsas {
						// fmt.Println(areaLsa)
						if areaLsa.LsaType == "routerLsa" {
							ospfRouter := OSPFRouter{
								RouterID: areaLsa.AdvertisingRouter,
							}
							for _, link := range areaLsa.OspfRouterLsa.RouterLsaLinks {
								if link.LinkType == "stubNetwork" {
									stub := Stub{
										Subnet: link.LinkData,
										Mask:   link.LinkID,
										Metric: link.Metric,
									}
									ospfRouter.Stubs = append(ospfRouter.Stubs, stub)
								} else if link.LinkType == "transitNetwork" {
									ospfRouter.Transits = append(ospfRouter.Transits, link.LinkData)
								} else if link.LinkType == "pointToPoint" {
									link := Link{
										Neighbor:  link.LinkID,
										IPAddress: link.LinkData,
										Metric:    link.Metric,
									}
									ospfRouter.Links = append(ospfRouter.Links, link)
								}
							}
							// fmt.Println(ospfRouter)
							ospfRoutersMap.Routers[ospfRouter.RouterID] = ospfRouter
							// OSPFRouters = append(OSPFRouters, ospfRouter)
						}
						if areaLsa.LsaType == "networkLsa" {
							networkLSADomain := areaLsa.LinkStateID
							networkLSASubnetMask := areaLsa.OspfNetworkLsa.NetworkMask
							cidr, err := getCIDR(networkLSADomain, networkLSASubnetMask)
							if err != nil {
								fmt.Println("Error getting CIDR: ", err)
								continue
							}
							// fmt.Println("CIDR: ", cidr)
							ospfRouter := OSPFRouter{
								RouterID: cidr,
							}
							for _, attachedRouter := range areaLsa.OspfNetworkLsa.AttachedRouters {
								neighborLink := Link{
									Neighbor:  cidr,
									IPAddress: "",
									Metric:    10,
								}
								neighborRouter := ospfRoutersMap.Routers[attachedRouter]
								for _, transit := range neighborRouter.Transits {
									isInSubnet, err := isIPAddressInSubnet(transit, cidr)
									if err != nil {
										fmt.Println("Error checking subnet: ", err)
										continue
									}
									if isInSubnet {
										neighborLink.IPAddress = transit
										break
									}
								}
								neighborRouter.Links = append(neighborRouter.Links, neighborLink)
								ospfRoutersMap.Routers[attachedRouter] = neighborRouter
							}
							ospfRoutersMap.Routers[ospfRouter.RouterID] = ospfRouter
							// OSPFRouters = append(OSPFRouters, ospfRouter)
						}
					}
				}
			}
		}
	}

	networkMap := &visualizer.NetworkMap{
		Devices: make(map[string][]visualizer.Edge),
	}

	for rid, ospfRouter := range ospfRoutersMap.Routers {
		networkMap.Devices[rid] = []visualizer.Edge{}

		for _, link := range ospfRouter.Links {
			edge := visualizer.Edge{
				LocalPort:    link.IPAddress,
				Neighbor:     link.Neighbor,
				NeighborPort: "",
			}

			neighborRID := link.Neighbor
			if !isIPAddressWithCIDR(neighborRID) {
				neighborObj := ospfRoutersMap.Routers[neighborRID]
				for _, neighborLink := range neighborObj.Links {
					if neighborLink.Neighbor == rid {
						edge.NeighborPort = neighborLink.IPAddress
						break
					}
				}
			}

			networkMap.Devices[rid] = append(networkMap.Devices[rid], edge)
		}
	}

	// Generate timestamp for filename (.json, .svg, .png)
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Exporting networkMap to a JSON file
	err = utils.SaveStructAsJson("ospf-graph", ospfRoutersMap, timestamp)
	if err != nil {
		fmt.Println("Error writing the file: ", err)
		return
	}

	// Exporting networkMap to a JSON file
	err = utils.SaveStructAsJson("graph", networkMap, timestamp)
	if err != nil {
		fmt.Println("Error writing the file: ", err)
		return
	}

	// Creating a graph and exporting it to a PNG file using go-graphviz
	err = visualizer.SaveTopologyWithGraphviz(ctx, networkMap, timestamp, "ospf")
	if err != nil {
		fmt.Println("Error rendering the topology: ", err)
		return
	}
}

// func getRouterByID(ospfRouters []OSPFRouter, routerID string) *OSPFRouter {
// 	for _, router := range ospfRouters {
// 		if router.RouterID == routerID {
// 			return &router
// 		}
// 	}
// 	return nil
// }

func getCIDR(ipAddress, subnetMask string) (string, error) {
	ip := net.ParseIP(ipAddress).To4()
	if ip == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	// Convert mask string to net.IPv4Mask
	maskParts := strings.Split(subnetMask, ".")
	if len(maskParts) != 4 {
		return "", fmt.Errorf("invalid subnet mask")
	}

	mask := net.IPv4Mask(
		parseByte(maskParts[0]),
		parseByte(maskParts[1]),
		parseByte(maskParts[2]),
		parseByte(maskParts[3]),
	)

	// Apply the mask to the IP to get the network address
	networkIP := ip.Mask(mask)

	// Get the prefix length
	ones, _ := mask.Size()

	cidr := fmt.Sprintf("%s/%d", networkIP.String(), ones)
	return cidr, nil
}

// Helper to convert string to byte (uint8)
func parseByte(s string) byte {
	var b int
	fmt.Sscanf(s, "%d", &b)
	return byte(b)
}

func isIPAddressInSubnet(ipStr string, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	_, subnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR: %s", cidrStr)
	}

	return subnet.Contains(ip), nil
}

func isIPAddressWithCIDR(input string) bool {
	_, _, err := net.ParseCIDR(input)
	return err == nil
}
