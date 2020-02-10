package packaging

import (
    "encoding/hex"
    "net"
    "syscall"

    "github.com/sirupsen/logrus"
)

// @localSocketConn is actual raw socket connection.
// which allow us to send the packet internally.
// @isSocketCreated represents whether socket variable
// is created or not.
var localSocketConn int
var isSocketCreated bool = false

// This method allow is to send the packet which is
// to be send internally to @dstIP and @dstPort.
func SendPacketInternally(dstIP net.IP,
			  dstPort uint16,
			  packetData []byte) error {

    // The raw socket had not created.
    if isSocketCreated == false {
	var err error
	localSocketConn, err = syscall.Socket(syscall.AF_INET,
					      syscall.SOCK_RAW,
					      syscall.IPPROTO_RAW)
	isSocketCreated = true
	if err != nil {
	    logrus.Debugf("Error during socket creation for" +
	                  "internal packet transmission \n")
	    return err
	}
    }
    // Creating senders Address.
    sendingAddress := syscall.SockaddrInet4{
	Port: int(dstPort),
    }
    // Appending IPv4 version of an address.
    copy(sendingAddress.Addr[:], dstIP.To4()[:4])

    // Sending packets internally.
    err := syscall.Sendto(localSocketConn,
			  packetData,
			  0,
			  &sendingAddress)
    if err != nil {
	logrus.Debugf("Error in transmission of packet to ipAddress %s and " +
	              "port %d \n", dstIP.String(), dstPort)
	logrus.Debugf("packet data is \n")
	logrus.Debugf(hex.Dump(packetData))
	return err
    }
    logrus.Debugf("packet send Internally")
    return nil
}