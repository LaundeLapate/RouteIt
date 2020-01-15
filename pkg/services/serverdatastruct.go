/*
This module is create to define all the dataType
that will be required and frequently used for the
Routing Server.
*/
package services

import "net"

// ClientsServer it is struct which combine various client's
// server parameters to clientID. So we can refer all the
// request related to clientServer can be handled by using
// clientID.
type ClientsServer struct {
	clientID            int16
	clientServerPort    int16
	clientServerNatPort int16
	// DeviceUsed is the DeviceName used by server's IP assigned
	//to client.
	DeviceUsed          string
	clientServerIP      net.IPAddr
	clientServerNatIP   net.IPAddr
}

// UpdateNatDetails allow us to update the NatDetails corresponding
// for the given ClientServer Addr.
func (p *ClientsServer) UpdateNatDetails(newNatIP net.IPAddr, newNatPort int16)  {
	p.clientServerNatIP   = newNatIP
	p.clientServerNatPort = newNatPort
}

// UpdateDeviceUsed allow us to update the Device used by the server's
// IP assigned to particular clientID.
func (p *ClientsServer) UpdateDeviceUsed(newDeviceName string)  err{
	p.DeviceUsed = newDeviceName
}

// ServerIPAssigned it is struct which combine clientID,
// to valid server IP and port as this server's IP and
// port will act as proxy IP address for client's server.
type ServerIPAssigned struct {
	assignedServerPort int16
	clientID           int16
	assignedServerIP   net.IPAddr
}

// FrameMapping is the struct which keep ARP mapping for all
// the devices as it keep mapping of all and internal mac Address.
type DeviceFrameMapping struct {
	DeviceName         string
	internalMacAddr    net.HardwareAddr
	externalMacAddr    net.HardwareAddr
}

// UpdateDeviceMacAddr Update's The macAddress as for a given Device.
func (p *DeviceFrameMapping) UpdateDeviceMacAddr(newInternalMacAddr net.HardwareAddr,
	                                             newExternalMacAddr net.HardwareAddr)  {
	p.internalMacAddr = newInternalMacAddr
	p.externalMacAddr = newExternalMacAddr
}
