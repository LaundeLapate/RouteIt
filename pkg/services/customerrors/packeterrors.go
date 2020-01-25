package customerrors

import "errors"

// @InProperCustomLayer
// This error is created is during parsing when ever size of custom layer is not correct.
// @ErrorInEthernetExtraction
// This is the error shows that Ethernet layer is not properly extracted.
// @ErrorInIPExtraction
// This is the error shows that IP layer is not properly extracted.
// @ErrorInTransportLayerExtraction
// This error show that we are not able to extract the transport layer.
// @WrongUDPCheckSum
// This Error is caused during computation of Checksum for the UDP layer.
// @WrongTCPCheckSum
// This Error is caused during computation of Checksum for the TCP layer.
// @LayerIsNotTCPOrUDP
// This Error show that above packet contains other than TCP or UDP layer as transport Layer.
var InProperCustomLayer             error = errors.New("size of custom layer is improper")
var ErrorInEthernetExtraction       error = errors.New("error in extracting ethernet layer")
var ErrorInIPExtraction             error = errors.New("packet is made of IPv4`")
var ErrorInTransportLayerExtraction error = errors.New("error in extracting IPv4 layer")
var WrongUDPCheckSum                error = errors.New("error in computing UDP checksum")
var WrongTCPCheckSum                error = errors.New("error in computing TCP checksum")
var LayerIsNotTCPOrUDP              error = errors.New("packet contains other than TCP or UDP layer")
var ARPtimeout                      error = errors.New("timeout during ARP request")