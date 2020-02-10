/*
This module consist of all the tools which is require for
extraction and updating information related Web application
cache.
 */
package webappcac

import "net"

// This function provides us the information about
// whether a the connection received from SrcIP to
// desIP is require to forward or not.
func ValidPacketConnection(srcIPAddr,
	                       dstIPAddr net.IPAddr,
	                       srcPort,
	                       dstPort   int16) bool {

}