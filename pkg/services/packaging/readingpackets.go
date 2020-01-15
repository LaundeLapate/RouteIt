/*
This module permits us to extract all the basic
details from the data packet as it permit us to
extract the information that is being extracted.
are:

isCustomPacket       Represents whether packet
					 contain external layer.
isPacketForPing      Provide information that
					 packet is have data or just
					 for ping.
<ethernetLayer>
<ipLayer>
<customLayer>        Layer which we will add for
				     internal communication
<remainingPayload>    [<transportLayer>   +
                       <applicationLayer> +
                       <actualPayload>](Unmodified)
 */
package packaging

import (
	"encoding/hex"
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// TransportInfo contains information regarding
// source and destination port of the the packet
// and transLayer from which further details can
// be extracted.
type TransportInfo struct {
	srcPort            int16
	dstPort 		   int16
	transLayer         gopacket.Layer
}

// This allow us extraction of transport from packet
// data and IPProtocol information.
func (p *TransportInfo) Init(packet *gopacket.Packet,
							 protocol layers.IPProtocol) error{

}

// PacketInfo is the struct which contains all
// the necessary about the packets of all kind
// non-custom, custom with data or ping.
type PacketInfo struct {
	isCustomPacket     bool
	isPacketForPing    bool
	ethernetLayer      *layers.Ethernet
	ipLayer            *layers.IPv4
	tspLayer		   TransportInfo
	customLayer        gopacket.Layer
	remainingPayload   gopacket.Payload
}

// This function initialise the PacketInfo from
// the packet. "packet" defines the actual packet
// which will be decoded, "hasCustomLayer" is
// the variable which provide information whether
// above packet has custom layer.
func (p *PacketInfo) ExtractInformation(packet gopacket.Packet,
										hasCustomLayer bool) error{

	p.isCustomPacket = hasCustomLayer

	ethernetLayerTmp := packet.Layer(layers.LayerTypeEthernet)
	// Checking whether ethernet layer is extracted properly.
	if ethernetLayerTmp == nil {
		errorManagement("EthernetLayer", "", packet)
		return errors.New("error in extracting ethernet layer")
	}
	p.ethernetLayer = ethernetLayerTmp.(*layers.Ethernet)

	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	// Checking whether IP layer is extracted properly.
	if ipv4Layer == nil {
		errorManagement("IPv4 layer", "", packet)
		return errors.New("error in extracting IPv4 layer")
	}
	p.ipLayer = ipv4Layer.(*layers.IPv4)

	// Extracting transport layer information.
	err := p.tspLayer.Init(&packet, p.ipLayer.Protocol)
	if err != nil {
		errorManagement("Transport layer",
						"IPProtocol was " + string(p.ipLayer.Protocol),
						packet)
		return errors.New("error in extracting IPv4 layer")
	}


	if p.isCustomPacket {
		// Extraction of custom layer.
	}

	return nil
}


// This function is created to throw error in
//formatted manner.
func errorManagement(errLayer, msg string, packet gopacket.Packet) {
	logrus.Debug("Error in extraction of " + errLayer + " of " +
		                  "packet")
	if len(msg) != 0 {
		logrus.Debug(msg)
	}
	// Dumping all the information in hexadecimal format.
	logrus.Debug("Pkt Information is", hex.Dump(packet.Data()))
}