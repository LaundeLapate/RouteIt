/*
This module consist of all the tools which is require for
extraction and updating information related routing server
cache.
 */
package servercac

import "github.com/LaundeLapate/RouteIt/pkg/services"


// This method provides us the ClientServer variable which
// is linked with client's ID.
func ProvideClientServerDetailsFromID(clientID int) services.ClientsServer {

}

// This method provides us ServerIpAssigned variable which
// is linked with client's ID.
func ProvideServerIPAssignedFromID(clientID int)  services.ServerIPAssigned {

}

// This method provide us the respective mac Address which are
// linked with Device name.
func DeviceDetailsFromDeviceName(deviceName string) services.DeviceFrameMapping {

}
