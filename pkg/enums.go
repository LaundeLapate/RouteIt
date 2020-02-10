package pkg

import (
    "net"

    "github.com/google/gopacket/layers"
)

// @EthernetInterface is the name of the ethernet interface of system.
// ServerAddr is the address provided by the NAT.
// @LocalInterface is the name of the local interface of system.
// @WireLessInterface is the name of the wireless interface of system.
// @IPLayerForPrototype	is a template for the IP layer
// @UDPLayerForPrototype is a template for the UDP layer.
var EthernetInterface    string     = "enx00e04c3607b7"
var ServerAddr           net.IPAddr = net.IPAddr{IP: net.ParseIP("127.0.0.1")}
var LocalInterface       string = "lo"
var WireLessInterface    string = "wlp2s0"
var emptyBaseLayer       layers.BaseLayer = layers.BaseLayer{Contents:[]byte{}, Payload:[]byte{}}
var IPLayerForPrototype  layers.IPv4      = layers.IPv4{BaseLayer:emptyBaseLayer, Version:4, IHL:5, TOS:0, Length:20, Id:9120, Flags:0, FragOffset:0, TTL:128, Protocol:17, Checksum:0, SrcIP:net.ParseIP("0.0.0.0"), DstIP:net.ParseIP("0.0.0.0")}
var UDPLayerForPrototype layers.UDP       = layers.UDP{BaseLayer:emptyBaseLayer, SrcPort:0, DstPort:0,Length:8, Checksum:0}

// @HolePunchPort is the port number provided for the hole punching.
// @IsPacketForPing is the value of byte to check whether packet is for pinging
// or not.
// @CustomLayerByteSize is lenght of custom layer.
const HolePunchPort       uint16 = 8000
const IsPacketForPing     uint8 = 1
const CustomLayerByteSize uint8 = 24

// This method allow us to intialise variable which
// will be required for further computation.
func init() {

}