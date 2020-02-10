/*
This module contain all the ket components which may be
requires during the creation of Ethernet layer.
 */
package packaging

import (
    "github.com/LaundeLapate/RouteIt/pkg/services/customerrors"
    "net"
    "strings"
    "time"

    "github.com/LaundeLapate/RouteIt/pkg"
    "github.com/google/gopacket/layers"
    "github.com/mdlayher/arp"
    "github.com/sirupsen/logrus"
)

// @AllConstructedClient represents the mapping of interface name
// to the address to arp.Client object.
// @AllConstructedInterfaces represents the mapping of interface name
// to the address to interface object.
// @AllInterfacesAddress represents the mapping of interface name
// to the address to local IP address.
// @ethernetClientCreate represents various various interface related
// variable are created or not.
// @timeOutCst timeout to all the ARP request.
var AllConstructedClients     = make(map[string]*arp.Client)
var AllConstructedInterfaces  = make(map[string]net.Interface)
var AllInterfacesAddress      = make(map[string]net.IP)
var EthernetType         int  = 0x0800
var ethernetClientCreate bool = false
var timeOutCst           time.Duration  = 1000 * time.Millisecond

// This function construct the ethernet frame corresponding
// to given @srcIPAdd and @dstIPAdd for given interfaceName.
// with ethernet type as IPv4
func GenerateEthernetLayer(interfaceName string,
			   srcIPAdd net.IP,
			   dstIPAdd net.IP) (layers.Ethernet, error) {

    var constructedEthernetLayer layers.Ethernet
    var err error

    // Checking for local interface.
    if interfaceName == pkg.LocalInterface {
        // Making Src and Dst mac Address as 00.00.00.00.00.00 .
	constructedEthernetLayer.DstMAC = net.HardwareAddr{00, 00, 00, 00, 00, 00}
	constructedEthernetLayer.SrcMAC = net.HardwareAddr{00, 00, 00, 00, 00, 00}
	constructedEthernetLayer.EthernetType = layers.EthernetType(EthernetType)
    }
    // Construction of ethernet interface.
    if interfaceName == pkg.EthernetInterface {
	// Checking whether apt.Client object corresponding
	// for given interface is created.
        if !ethernetClientCreate {
	    err = createClient(interfaceName)
	    if err != nil {
		logrus.Debug("Can't create a client variable " +
			      "for %s", interfaceName)
		return constructedEthernetLayer, err
	    }
	    // setting that arp.Client object is created for given
	    // interface.
	    ethernetClientCreate = true
	}

	// Extraction arp.Client object.
	clientName := AllConstructedClients[interfaceName]
	// Resolving the ethernet address for given
	// IP address and interface.
	constructedEthernetLayer.SrcMAC, err = ResolveWrapper(*clientName,
							      srcIPAdd,
							      interfaceName)
	constructedEthernetLayer.DstMAC, err = ResolveWrapper(*clientName,
	                                                      dstIPAdd,
							      interfaceName)
	constructedEthernetLayer.EthernetType = layers.EthernetType(EthernetType)
    }

    // Construction of wired interface.
    if interfaceName == pkg.WireLessInterface {

    }
    return constructedEthernetLayer, nil
}

// This method create the apr.Client object corresponding
// to given interface name.
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

    // Mapping interface to arp.Client object.
    AllConstructedClients[interfaceName] = clientForInterface
    return errForClient
}

// This function is just wrapper above the actual resolve function
// which allow us to know the mac address corresponding to IP address
// and interface. This wrapper just have additional feature of timout
func ResolveWrapper(client arp.Client,
                    ipAddr net.IP,
                    interfaceName string) (net.HardwareAddr, error) {

    // constructing the IP Address of device
    if AllInterfacesAddress[interfaceName].String() == ipAddr.String() {
		return client.HardwareAddr(), nil
	}

    // Creating the channel for timeout.
    timeOutChannel := make(chan string, 1)
    var macAddr net.HardwareAddr
    var err error

    // Extracting the mac address for given IP and interface.
    go func() {
	macAddr, err = client.Resolve(ipAddr)
	if err != nil {
	    logrus.Debug("Error in Resolving mac Address in %s and " +
	                 "IP address", client.HardwareAddr(), ipAddr.String())
	    logrus.Debug(err)
	}

    // passing message when mac Address is computed.
    timeOutChannel <- "computed mac address"
    }()

    select {
    // When mac Address is computed.
    case <- timeOutChannel:
	if err != nil {
			return macAddr, err
		}
	    return macAddr, nil
    // Timout trigger.
    case <- time.After(timeOutCst):
	logrus.Debug("Mac Address of %s and %d can't be computed due " +
		     "to timeout", client.HardwareAddr(), ipAddr.String())
	return macAddr, customerrors.ARPtimeout
    }
}
