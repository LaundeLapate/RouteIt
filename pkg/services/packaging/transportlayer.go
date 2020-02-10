/*
This module has Data Structures and method
which will be required to extract information
from the packet.
 */
package packaging

import (
    "github.com/LaundeLapate/RouteIt/pkg/services/customerrors"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
)

// TransportInfo contains information regarding
// source and destination port of the packet and
// transLayer from which further details can
// be extracted.
type TransportInfo struct {
	SrcPort            uint16
	DstPort 	   uint16
	TransLayer         gopacket.Layer
}

// This allows us extraction of transport from packet
// data and IPProtocol information.
func (p *TransportInfo) Init(packet *gopacket.Packet,
                             protocol layers.IPProtocol) error {

    // Extracting various parameters on the basic of
    // type of transport layer protocol.
    switch protocol {
    // layer is TCP.
    case layers.IPProtocolTCP:
	p.TransLayer = (*packet).Layer(layers.LayerTypeTCP)
	tmpLayerValue := p.TransLayer.(*layers.TCP)
	p.SrcPort = uint16(tmpLayerValue.SrcPort)
	p.DstPort = uint16(tmpLayerValue.DstPort)

    // layer is UDP
    case layers.IPProtocolUDP:
	p.TransLayer = (*packet).Layer(layers.LayerTypeUDP)
	tmpLayerValue := p.TransLayer.(*layers.UDP)
	p.SrcPort = uint16(tmpLayerValue.SrcPort)
	p.DstPort = uint16(tmpLayerValue.DstPort)
    default:
	return customerrors.LayerIsNotTCPOrUDP
    }
    return nil
}
