package pkg

import "net"

// @EthernetInterface is name of the ethernet interface of system.
// ServerAddr is the address provided by the NAT.
// @LocalInterface is the name of the local interface of system.
// @WireLessInterface is the name of the wireless interface of system.
var EthernetInterface   string     = "enx00e04c3607b7"
var ServerAddr          net.IPAddr = net.IPAddr{IP: net.ParseIP("10.38.2.35")}
var LocalInterface    string = "lo"
var WireLessInterface string = "wlp2s0"

// @HolePunchPort is the port number provided for hole punching.
// @IsPacketForPing is the value of byte to check whether packet is for pinging
// or not.
// @CustomLayerByteSize is lenght of custom layer.
const HolePunchPort       uint16 = 8000
const IsPacketForPing     uint8 = 1
const CustomLayerByteSize uint8 = 17


// This method allow us to intialise variable which
// will be required for further computation.
func init() {

}