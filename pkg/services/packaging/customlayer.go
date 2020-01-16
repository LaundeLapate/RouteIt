/*
This module manage the custom layer that will be
that is being added with the packet for the internal
communication purpose.
*/
package packaging

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/google/gopacket"
)

// Creating CustomLayerType layer variable
// for our custom layer.
// And linking Custom layer type to our Decode
// function. i.e "DecodeCstLayer"
var CustomLayerType = gopacket.RegisterLayerType(
	2001,
	gopacket.LayerTypeMetadata{
		Name:    "CustomLayer",
		Decoder: gopacket.DecodeFunc(DecodeCstLayer),
	},
)

// Creating custom data type for Custom layer.
type CustomLayer struct {
	// IsPing represents whether the packet is for
	// pinging or has actual data tobe transported.
	IsPing         uint8
	ClientSeverID  uint64
	ClientIP       net.IPAddr
	ClientPort     uint16
	// Actual payload of packet.
	restOfData     []byte
}

// DecodeFromBytes extract all the attributes for
// CustomLayer from the byte data.
func (l *CustomLayer) DecodeFromBytes(data []byte) error {
	// Checks whether dataPacket size is less than
	// 16 bytes as layer data can't be less than.
	if len(data) < 16 {
		return errors.New("header Size is less than 16 byte")
	}

	l.IsPing = data[0]
	l.ClientSeverID = binary.BigEndian.Uint64(data[1:9])
	l.ClientPort    = binary.BigEndian.Uint16(data[10:12])
	l.ClientIP = net.IPAddr{IP: data[12:16]}
	l.restOfData = data[16:]

	return nil
}

// DecodeCustomLayer merges all the CustomLayer's
// parameter into slice of byte which is 16 byte long.
func (l *CustomLayer) DecodeCustomLayer() []byte {
	layerData := make([]byte, 16)
	layerData[0] = l.IsPing
	binary.BigEndian.PutUint64(layerData[1:9], l.ClientSeverID)
	binary.BigEndian.PutUint16(layerData[10:12], l.ClientPort)
	copy(layerData[12:16], l.ClientIP.IP)
	return layerData
}

// Integrating with Layer Interface too.
// LayerType provide information of layer.
func (l *CustomLayer) LayerType() gopacket.LayerType {
	return CustomLayerType
}

// Provides content of layer in bytes.
func (l *CustomLayer) LayerContents() []byte  {
	return l.DecodeCustomLayer()
}

// Provides payload of our packet.
func (l *CustomLayer) LayerPayload() []byte {
	return l.restOfData
}

// DecodeCstLayer function integrate's CustomLayer
// to gopacket modules. As here we are extraction
// custom layer form the data and initializing
// payLoad as NewEncoding point.
func DecodeCstLayer(data []byte, p gopacket.PacketBuilder) error {
	// Initializing CustomLayer from data.
	customLayer := &CustomLayer{}
	err := customLayer.DecodeFromBytes(data)
	if err != nil {
		return errors.New("custom layer can't be extracted")
	}

	// Adding customLayer to our packet.
	p.AddLayer(customLayer)
	// Initializing payLoad as LayerTypePayload.
	return p.NextDecoder(gopacket.LayerTypePayload)
}
