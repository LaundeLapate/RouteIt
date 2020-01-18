/*
This module permits us to create link layer for
given request.
 */
package packaging

import (
	"errors"
	"net"
	"time"

	"github.com/LaundeLapate/RouteIt/pkg"
	"github.com/google/gopacket/layers"
	"github.com/mdlayher/arp"
	"github.com/sirupsen/logrus"
)

var AllConstructedClients map[string]*arp.Client
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
			AllConstructedClients[interfaceName], err = createClient(interfaceName)
			if err != nil {
				logrus.Debug("Can't create a client variable " +
					                  "for %s", interfaceName)
				return constructedEthernetLayer, err
			}
			ethernetClientCreate = true
		}
		clientName := AllConstructedClients[interfaceName]
		constructedEthernetLayer.SrcMAC, err  = ResolveWrapper(*clientName, srcIPAdd)
		constructedEthernetLayer.DstMAC, err  = ResolveWrapper(*clientName, dstIPAdd)
		constructedEthernetLayer.EthernetType = 0x0800
	}

	// Construction of wired interface.
	if interfaceName == pkg.WireLessInterface {

	}

	return constructedEthernetLayer, nil
}

func createClient(interfaceName string) (*arp.Client, error){
	interfaceVar, err := net.InterfaceByName(interfaceName)
	if err != nil {
		logrus.Debug("interface for %s can't be determined", interfaceName)
		return nil, err
	}

	clientForInterface, errForClient := arp.Dial(interfaceVar)
	if errForClient != nil {
		logrus.Debug("client't can't be created for %s interface", clientForInterface)
	}
	AllConstructedClients[interfaceName] = clientForInterface
}

// Wrapper over resolve which allow us to macAddress
// corresponding to device along with a timeout.
func ResolveWrapper(client arp.Client, ipAddr net.IP) (net.HardwareAddr, error) {
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
		logrus.Debug("Mac Address of %s and %d can't be computed",
							  client.HardwareAddr(), ipAddr.String())
		return macAddr, errors.New("timout mac address can't be computed")
	}
}