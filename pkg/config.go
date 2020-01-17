package pkg

import "net"

const ServerAddr      string = "127.0.0.1"
const HolePunchPort   uint16 = 8000
const IsPacketForPing  uint8 = 1


var EthernetInterface   string
var ServerHolePunchAddr net.IPAddr = net.IPAddr{IP: net.ParseIP(ServerAddr)}
