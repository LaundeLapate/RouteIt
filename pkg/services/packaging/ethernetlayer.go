/*
This module permits us to create link layer for
given request.
 */
package packaging

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/LaundeLapate/RouteIt/pkg"
	"github.com/google/gopacket/layers"
	"github.com/mdlayher/arp"
	"github.com/sirupsen/logrus"
)

var AllConstructedClients    = make(map[string]*arp.Client)
var AllConstructedInterfaces = make(map[string]net.Interface)
var AllInterfacesAddress     = make(map[string]net.IP)
var ethernetClientCreate bool = false
var timeOutCst           time.Duration  = 1000 * time.Millisecond

// This function construct the ethernet layer for
// given interfaceName.
func GenerateEthernetLayer(interfaceName string,
	                       srcIPAdd net.IP,
	                       dstIPAdd net.IP) (layers.Ethernet, error) {

	var constructedEthernetLayer layers.Ethernet
	var err error
	// Checking for local interface.
	if interfaceName == pkg.LocalInterface {
		constructedEthernetLayer.DstMAC = net.HardwareAddr{00, 00, 00, 00, 00, 00}
		constructedEthernetLayer.SrcMAC = net.HardwareAddr{00, 00, 00, 00, 00, 00}
		constructedEthernetLayer.EthernetType = 0x0800
	}

	// Construction of ethernetFrame
	if interfaceName == pkg.EthernetInterface {
		if !ethernetClientCreate {
			err = createClient(interfaceName)
			if err != nil {
				logrus.Debug("Can't create a client variable " +
					                  "for %s", interfaceName)
				return constructedEthernetLayer, err
			}
			ethernetClientCreate = true
		}
		clientName := AllConstructedClients[interfaceName]
		constructedEthernetLayer.SrcMAC, err  = ResolveWrapper(*clientName,
																srcIPAdd,
																interfaceName)
		constructedEthernetLayer.DstMAC, err  = ResolveWrapper(*clientName,
																dstIPAdd,
																interfaceName)
		constructedEthernetLayer.EthernetType = 0x0800
	}

	// Construction of wired interface.
	if interfaceName == pkg.WireLessInterface {

	}

	return constructedEthernetLayer, nil
}

func createClient(interfaceName string) error {
	interfaceVar, err := net.InterfaceByName(interfaceName)
	if err != nil {
		logrus.Debug("interface for %s can't be determined", interfaceName)
		return  err
	}
	// Constructing interface for given device.
	AllConstructedInterfaces[interfaceName] = *interfaceVar

	// Extracting address for given device.
	allAddress, _ := (*interfaceVar).Addrs()
	addressInString := strings.Split(allAddress[0].String(), "/")[0]
	AllInterfacesAddress[interfaceName] = net.ParseIP(addressInString)

	clientForInterface, errForClient := arp.Dial(interfaceVar)
	if errForClient != nil {
		logrus.Debug("client't can't be created for %s " +
							  "interface", clientForInterface)
	}

	AllConstructedClients[interfaceName] = clientForInterface
	return errForClient
}

// Wrapper over resolve which allow us to macAddress
// corresponding to device along with a timeout.
func ResolveWrapper(client arp.Client, ipAddr net.IP, interfaceName string) (net.HardwareAddr, error) {
	// constructing the IP Address of device
	if AllInterfacesAddress[interfaceName].String() == ipAddr.String() {
		return client.HardwareAddr(), nil
	}

	timeOutChannel := make(chan string, 1)
	var macAddr net.HardwareAddr
	var err error
	go func() {
		macAddr, err = client.Resolve(ipAddr)
		if err != nil {
			logrus.Debug("Error in Resolving mac Address in %s and " +
				                  "IP address", client.HardwareAddr(), ipAddr.String())
			logrus.Debug(err)
		}
		timeOutChannel <- "computer mac address"
	}()

	select {
	case <- timeOutChannel:
		if err != nil {
			return macAddr, err
		}
		return macAddr, nil
	case <- time.After(timeOutCst):
		logrus.Debug("Mac Address of %s and %d can't be computed due " +
							"to timeout", client.HardwareAddr(), ipAddr.String())
		return macAddr, errors.New("timout mac address can't be computed")
	}
}