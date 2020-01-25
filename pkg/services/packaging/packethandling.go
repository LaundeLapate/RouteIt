/*
This module permits us to extract all the basic
details from the data packet as it permit us to
extract the information that is being extracted.
are:
 */
package packaging

import (
    "github.com/LaundeLapate/RouteIt/pkg/services/customerrors"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/sirupsen/logrus"
)

// PacketInfo is the struct which contains all
// the necessary about the packets of all kind
// non-custom, custom with data or ping.
type PacketInfo struct {
    EthernetLayer      layers.Ethernet
    IpLayer            layers.IPv4
    TspLayer           TransportInfo
    AdditionalLayer    CustomLayer
    RemainingPayload   gopacket.Payload
}

// This function initialise the PacketInfo from
// the packet. "packet" defines the actual packet
// which will be decoded, "hasCustomLayer" is
// the variable which provide information whether
// above packet has custom layer.
func (p *PacketInfo) ExtractInformation(packet gopacket.Packet,
                                        linkLayerIsThere bool) error {

    // Extracting link layer for ethernet device.
    if linkLayerIsThere {
	extractedEthernetLayer := packet.Layer(layers.LayerTypeEthernet)

	// Checking whether ethernet layer is extracted properly.
	if extractedEthernetLayer == nil {
	    logrus.Debugf("Error during extraction of linkLayer \n")
	    return customerrors.ErrorInEthernetExtraction
	}
	p.EthernetLayer = *(extractedEthernetLayer.(*layers.Ethernet))
    }

    // Extracting IPv4 layer.
    ipv4Layer := packet.Layer(layers.LayerTypeIPv4)

    // Checking whether IP layer is extracted properly.
    if ipv4Layer == nil {
	logrus.Debugf("IPv6 data packet \n")
	return customerrors.ErrorInIPExtraction
    }
    p.IpLayer = *(ipv4Layer.(*layers.IPv4))

    // Extracting transport layer information.
    err := p.TspLayer.Init(&packet, p.IpLayer.Protocol)
    if err != nil {
        logrus.Debugf("Error in Extraction in transport layer \n")
        logrus.Debug(err)
        return customerrors.ErrorInTransportLayerExtraction
    }

    if packet.ApplicationLayer() != nil {
	    p.RemainingPayload = packet.ApplicationLayer().Payload()
    }

    return nil
}

// This method construct the packets data in to bytes
// which we will able to send throw wires.
func (p *PacketInfo) ConstructPacket(buffer  *gopacket.SerializeBuffer,
                                     options *gopacket.SerializeOptions,
                                     sendInternally bool,
                                     interfaceName string) ([]byte, error) {

    // Appending custom layer at start of payLoad.
    customPayLoad := p.RemainingPayload
    // Creating slice to keep all the layers.
    var allLayers []gopacket.SerializableLayer
    var err error

    // Options shows various parameter that that is
    // lenght of new packet is same and we have to
    // recompute the checksum.
    options = &gopacket.SerializeOptions{FixLengths: false,
                                         ComputeChecksums: true}

    // Checking whether packet should have link layer.
    if sendInternally {
	// Adding linkLayer as packet must
	// contain linkLayer.
	p.EthernetLayer, err = GenerateEthernetLayer(interfaceName,
						     p.IpLayer.SrcIP,
						     p.IpLayer.DstIP)

	// Checking for error in link layer.
	if err != nil {
	    logrus.Debug("Error in creating the layer link layer")
	    return []byte{}, err
	}
	allLayers = append(allLayers, &p.EthernetLayer)
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
		return []byte{}, customerrors.WrongUDPCheckSum
	}
	allLayers = append(allLayers, udpLayerToBeAdded)

    case layers.IPProtocolTCP:
	// Addition of TCP layer.
	tcpLayerToBeAdded := p.TspLayer.TransLayer.(*layers.TCP)
	tcpLayerToBeAdded.DstPort = layers.TCPPort(p.TspLayer.DstPort)
	tcpLayerToBeAdded.SrcPort = layers.TCPPort(p.TspLayer.SrcPort)
	err := tcpLayerToBeAdded.SetNetworkLayerForChecksum(&p.IpLayer)
	if err != nil {
		return []byte{}, customerrors.WrongTCPCheckSum
	}
	allLayers = append(allLayers, tcpLayerToBeAdded)

    default:
	return []byte{}, customerrors.LayerIsNotTCPOrUDP
    }

    // Adding payload to the packet.
    allLayers = append(allLayers, customPayLoad)

    // Serialising the buffer.
    err = gopacket.SerializeLayers(*buffer, *options, allLayers...)

    if err != nil {
	    return []byte{}, nil
    }
    return (*buffer).Bytes(), nil
}
