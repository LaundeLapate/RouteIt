/*
This module permits us to extract all the basic
details from the data packet as it permit us to
extract the information that is being extracted.
are:

// TODO: rewrite the definition.
// TODO: Add separate error variable for IPv6 and
// TODO: and none TCP and UDP data packets.
// TODO: solution it data packet size exceed packet size.
 */
package packaging

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// PacketInfo is the struct which contains all
// the necessary about the packets of all kind
// non-custom, custom with data or ping.
type PacketInfo struct {
	IsCustomPacket     bool
	IsPacketForPing    bool
	EthernetLayer      layers.Ethernet
	IpLayer            layers.IPv4
	TspLayer		   TransportInfo
	AdditionalLayer    CustomLayer
	RemainingPayload   gopacket.Payload
}

// This function initialise the PacketInfo from
// the packet. "packet" defines the actual packet
// which will be decoded, "hasCustomLayer" is
// the variable which provide information whether
// above packet has custom layer.
func (p *PacketInfo) ExtractInformation(packet gopacket.Packet,
										hasCustomLayer bool) error{

	p.IsCustomPacket = hasCustomLayer

	ethernetLayerTmp := packet.Layer(layers.LayerTypeEthernet)
	// Checking whether ethernet layer is extracted properly.
	if ethernetLayerTmp == nil {
		errorManagement("EthernetLayer", "", packet)
		return errors.New("error in extracting ethernet layer")
	}
	p.EthernetLayer = *(ethernetLayerTmp.(*layers.Ethernet))

	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	// Checking whether IP layer is extracted properly.
	if ipv4Layer == nil {
		errorManagement("IPv4 layer", "", packet)
		return errors.New("IPv6 is being used`")
	}
	p.IpLayer = *(ipv4Layer.(*layers.IPv4))

	// Extracting transport layer information.
	err := p.TspLayer.Init(&packet, p.IpLayer.Protocol)
	if err != nil {
		errorManagement("Transport layer",
						"IPProtocol was " + string(p.IpLayer.Protocol),
						packet)
		return errors.New("error in extracting IPv4 layer")
	}

	if packet.ApplicationLayer() != nil {
		p.RemainingPayload = packet.ApplicationLayer().Payload()
	}

	if p.IsCustomPacket {
		err := p.AdditionalLayer.DecodeFromPayload(p.RemainingPayload)
		if err != nil {
			return err
		}
		p.RemainingPayload = p.AdditionalLayer.TspLayer.TransLayer.LayerPayload()
	}
	return nil
}

// This method permit us to construct packet of both
// type which can be send throw wire.
// buffer represents the actual buffer to which packet
// is going to be transmitted.
// addEthernet tell us to whether to add ethernet layer
//or not.
func (p *PacketInfo) ConstructPacket(buffer *gopacket.SerializeBuffer,
									 options *gopacket.SerializeOptions,
									 addEthernet bool) ([]byte, error) {

	// Appending custom layer at start of payLoad.
	//customPayLoad := gopacket.Payload(append(p.RemainingPayload.LayerContents(),
	//								         p.RemainingPayload...))
	customPayLoad := p.RemainingPayload
	// Creating slice to keep all the layers.
	var allLayers []gopacket.SerializableLayer

	options = &gopacket.SerializeOptions{FixLengths:       true,
										 ComputeChecksums: true}
	// Adding ethernet layer as packet must
	// contain ethernet layer.
	if addEthernet {
		allLayers = append(allLayers, &p.EthernetLayer)
	} else {
		newEthernetLayer := layers.Ethernet{SrcMAC: net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			                               DstMAC: net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			                               EthernetType: 0x0800}
		fmt.Println(hex.Dump(newEthernetLayer.LayerContents()))
		allLayers = append(allLayers, &newEthernetLayer)
	}
	// Adding IP layer.
	allLayers = append(allLayers, &p.IpLayer)
	// Serializing buffer on the basis of
	// type of transport layer.
	*buffer = gopacket.NewSerializeBuffer()
	switch p.IpLayer.Protocol {
	case layers.IPProtocolUDP:
		// Addition of UDP layer.
		udpLayerToBeAdded := p.TspLayer.TransLayer.(*layers.UDP)
		udpLayerToBeAdded.SrcPort = layers.UDPPort(p.TspLayer.SrcPort)
		udpLayerToBeAdded.DstPort = layers.UDPPort(p.TspLayer.DstPort)
		err := udpLayerToBeAdded.SetNetworkLayerForChecksum(&p.IpLayer)
		if err != nil {
			return []byte{}, errors.New("error in handling udp check sum")
		}
		allLayers = append(allLayers, udpLayerToBeAdded)

	case layers.IPProtocolTCP:
		// Addition of TCP layer.
		tcpLayerToBeAdded := p.TspLayer.TransLayer.(*layers.TCP)
		tcpLayerToBeAdded.DstPort = layers.TCPPort(p.TspLayer.DstPort)
		tcpLayerToBeAdded.SrcPort = layers.TCPPort(p.TspLayer.SrcPort)
		err := tcpLayerToBeAdded.SetNetworkLayerForChecksum(&p.IpLayer)
		if err != nil {
			return []byte{}, errors.New("error in handling tcp check sum")
		}
		allLayers = append(allLayers, tcpLayerToBeAdded)

	default:
		return []byte{}, errors.New("error in handling tcp check sum")

	}

	// Adding payload to the packet.
	allLayers = append(allLayers, customPayLoad)

	// Serialising the buffer.
	err := gopacket.SerializeLayers(*buffer, *options, allLayers...)

	if err != nil {
		return []byte{}, nil
	}
	return (*buffer).Bytes(), nil
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