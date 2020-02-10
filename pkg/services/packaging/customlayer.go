/*
This module has basic dataStructures and functions
that may be required for creation and management of
custom layer.
*/
package packaging

import (
	"encoding/binary"
	"net"

	"github.com/LaundeLapate/RouteIt/pkg"
	"github.com/LaundeLapate/RouteIt/pkg/services/customerrors"
	"github.com/sirupsen/logrus"
)

type CustomLayer struct {
    // IsPing represents whether the packet is for
    // pinging or has actual data tobe transported.
    IsPing        uint8
    ClientSeverID uint64
    ClientIP      net.IP
    ClientPort    uint16
}

// This method allow us to convert CustomLayer struct
// to its respective byte data that is of fix length
// @CustomLayerByteSize.
func (l *CustomLayer) CovertCustomLayerToBytes() []byte {
    var convertedData = make([]byte, pkg.CustomLayerByteSize)
    convertedData[0] = l.IsPing

    // Adding clientID to the custom layer.
    binary.BigEndian.PutUint64(convertedData[1:9], l.ClientSeverID)

    // Appending clientIP and Port to the custom
    // layer.
    binary.BigEndian.PutUint16(convertedData[10:12], l.ClientPort)
    copy(convertedData[12:16], l.ClientIP)

    return convertedData
}

// This method all us to intialise the custom layer from
// provided byte data.
func (l *CustomLayer) CreateLayerFromByte(customLayerData []byte) error {

    // Validating whether layer is proper or not.
    if len(customLayerData) != int(pkg.CustomLayerByteSize) {
	logrus.Debugf("Unable to extract custom layer as improper size. \n")
	return customerrors.InProperCustomLayer
    }
    // Extracting various parameters.
    l.IsPing        = customLayerData[0]
    l.ClientSeverID = binary.BigEndian.Uint64(customLayerData[1:9])
    l.ClientPort    = binary.BigEndian.Uint16(customLayerData[10:12])
    l.ClientIP      = customLayerData[12:16]

    return nil
}
