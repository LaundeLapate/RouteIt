/*
This module have Data Structures and method
which will be require to extract information
from packet.
 */
package packaging

import (
	"errors"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TransportInfo contains information regarding
// source and destination port of the the packet
// and transLayer from which further details can
// be extracted.
type TransportInfo struct {
	SrcPort            uint16
	DstPort 		   uint16
	TransLayer         gopacket.Layer
}

// This allow us extraction of transport from packet
// data and IPProtocol information.
func (p *TransportInfo) Init(packet *gopacket.Packet,
	                         protocol layers.IPProtocol) error{

	switch protocol {
	case layers.IPProtocolTCP:
		p.TransLayer = (*packet).Layer(layers.LayerTypeTCP)
		tmpLayerValue := p.TransLayer.(*layers.TCP)
		p.SrcPort = uint16(tmpLayerValue.SrcPort)
		p.DstPort = uint16(tmpLayerValue.DstPort)

	case layers.IPProtocolUDP:
		p.TransLayer = (*packet).Layer(layers.LayerTypeUDP)
		tmpLayerValue := p.TransLayer.(*layers.UDP)
		p.SrcPort = uint16(tmpLayerValue.SrcPort)
		p.DstPort = uint16(tmpLayerValue.DstPort)
	default:
		return errors.New("packet has protocols other than" +
			                    "TCP and UDP")
	}
	return nil
}

func (p *TransportInfo) CreateTspLayerFromByte(data []byte,
	                                           protocol layers.IPProtocol) error {

	switch protocol {
	case layers.IPProtocolTCP:
		tmpTCPVar := layers.TCP{}
		var tmp gopacket.DecodeFeedback
		err := tmpTCPVar.DecodeFromBytes(data, tmp)
		if err != nil {
			return err
		}

		p.TransLayer = &tmpTCPVar
		p.SrcPort = uint16(tmpTCPVar.SrcPort)
		p.DstPort = uint16(tmpTCPVar.DstPort)

	case layers.IPProtocolUDP:
		tmpUDPVar := layers.UDP{}
		var tmp gopacket.DecodeFeedback
		err := tmpUDPVar.DecodeFromBytes(data, tmp)
		if err != nil {
			return err
		}

		p.TransLayer = &tmpUDPVar
		p.SrcPort = uint16(tmpUDPVar.SrcPort)
		p.DstPort = uint16(tmpUDPVar.DstPort)

	default:
		return errors.New("packet has protocols other than" +
			                    "TCP and UDP")
	}

	return nil
}
