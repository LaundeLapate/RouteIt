/*
This module provides functionality of adding
and combining and parsing custom layer with
an actual packet.
 */

package packaging

import (
    "github.com/google/gopacket"
    "net"

    "github.com/LaundeLapate/RouteIt/pkg"
    "github.com/google/gopacket/layers"
    "github.com/sirupsen/logrus"
)

// This method allow us to add combine actual
// packet with custom layer to make is
// transportable from the punched hole over the
// NAT.
func AddCustomLayerToPacketInfo(dstIP   net.IP,
    				dstPort uint16,
                                nonCustomPacket PacketInfo,
                                customLayerParameters CustomLayer) PacketInfo {

    customPacket := PacketInfo{}
    newIPLayer   := pkg.IPLayerForPrototype
    newUDPLayer  := pkg.UDPLayerForPrototype

    // Constructing IP layer for a packet.
    newIPLayer.DstIP    = dstIP
    newIPLayer.SrcIP    = pkg.ServerAddr.IP
    newIPLayer.Protocol = layers.IPProtocolUDP

    // Creating the custom payload as combination
    // custom layer and all the value from above
    // Ethernet layer.
    payloadFromIPlayer := append(nonCustomPacket.IpLayer.Contents, nonCustomPacket.IpLayer.Payload...)
    payloadForCustomLayer := append(customLayerParameters.CovertCustomLayerToBytes(),
                                    payloadFromIPlayer...)

    // Constructing Transport layer for the packet.
    newUDPLayer.SrcPort           = layers.UDPPort(pkg.HolePunchPort)
    customPacket.TspLayer.SrcPort = pkg.HolePunchPort
    newUDPLayer.DstPort           = layers.UDPPort(dstPort)
    customPacket.TspLayer.DstPort = dstPort
    newUDPLayer.Payload = payloadForCustomLayer

    // Updating the lenght of various layers.
    newUDPLayer.Length += uint16(len(payloadForCustomLayer))
    newIPLayer.Length  += newUDPLayer.Length

    // Appending all the parameters to our
    // custom layer.
    customPacket.IpLayer = newIPLayer
    customPacket.TspLayer.TransLayer = &newUDPLayer
    customPacket.RemainingPayload = payloadForCustomLayer

    return customPacket
}

// This method allow us to extract custom layer
// from the provided custom packets.
func ExtractCustomLayer(customPacket PacketInfo) (PacketInfo, CustomLayer, error) {

    // custom layer information in bytes.
    customLayerInformation := customPacket.RemainingPayload[:pkg.CustomLayerByteSize]

    // Removing bytes corresponding to custom
    // layer.
    customPacket.RemainingPayload = customPacket.RemainingPayload[pkg.CustomLayerByteSize:]

    nonCustomPacketData := append(customPacket.EthernetLayer.Contents,
                                  customPacket.RemainingPayload...)

    // Creating new packet corresponding to the bytes set.
    newPacketCreated := gopacket.NewPacket(nonCustomPacketData,
                                           layers.LayerTypeEthernet,
                                           gopacket.Default)

    ExtractedPacketInfo := &PacketInfo{}

    // Extracting information from the custom layer.
    customPacketInfo := CustomLayer{}
    err := customPacketInfo.CreateLayerFromByte(customLayerInformation)
    if err != nil {
        logrus.Debugf("Error during extraction information for the " +
                      "for the custom struct")
        return *ExtractedPacketInfo, customPacketInfo, err
    }

    // Extracting packet information from the byte's data.
    err = ExtractedPacketInfo.ExtractInformation(newPacketCreated,
                                true)

    if err != nil {
        logrus.Debugf("Error during extraction of information " +
                      "from newly created from the custom Packet \n")
        return *ExtractedPacketInfo, customPacketInfo, err
    }

    return *ExtractedPacketInfo, customPacketInfo, nil
}
