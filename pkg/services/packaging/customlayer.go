package packaging

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/LaundeLapate/RouteIt/pkg"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type CustomLayer struct {
	// IsPing represents whether the packet is for
	// pinging or has actual data tobe transported.
	IsPing        uint8
	ClientSeverID uint64
	ClientIP      net.IPAddr
	ClientPort    uint16
	IPProtocol    layers.IPProtocol
	// Transport layer of old packet.
	TspLayer      TransportInfo
}

// This method allow us to extract custom layer
// from the payload from where it has been extracted.
func (l *CustomLayer) DecodeFromPayload(payload gopacket.Payload) error {
	var data = []byte(payload)
	// Checking whether packet is at least for ping
	// message.
	if len(data) < 9 {
		return errors.New("custom layer is not create properly")
	}

	l.IsPing = data[0]
	l.ClientSeverID = binary.BigEndian.Uint64(data[1:9])

	// Packet was there for ping message.
	if l.IsPing == pkg.IsPacketForPing {
		return nil
	}

	// Since packet is not for ping message therefore it
	// must have least length of 17.
	if len(data) < 17 {
		return errors.New("header Size is less than 16 byte")
	}

	l.ClientPort    = binary.BigEndian.Uint16(data[10:12])
	l.ClientIP   = net.IPAddr{IP: data[12:16]}
	l.IPProtocol = layers.IPProtocol(int8(data[16]))
	data = data[16:]

	// Creating Transport Layer from given data bits.
	err := l.TspLayer.CreateTspLayerFromByte(data, l.IPProtocol)
	if err != nil {
		return err
	}
	return nil
}

// This layer allow us to encode the CustomLayer
// struct into. byte data which can be further
// merged with payload.
func (l *CustomLayer) EncodeTransportInfo() []byte {
	layerData := make([]byte, 9)
	if l.IsPing != pkg.IsPacketForPing {
		layerData = make([]byte, 17)
	}
	layerData[0] = l.IsPing
	binary.BigEndian.PutUint64(layerData[1:9], l.ClientSeverID)

	// Checking whether custom layer is for
	// pinging purpose if so we only send
	// clientID.
	if l.IsPing == pkg.IsPacketForPing {
		return layerData
	}

	// Appending remaining data in other case.
	binary.BigEndian.PutUint16(layerData[10:12], l.ClientPort)
	copy(layerData[12:16], l.ClientIP.IP)
	layerData[16] = byte(l.IPProtocol)
	layerData = append(layerData, l.TspLayer.TransLayer.LayerContents()...)
	return layerData
}