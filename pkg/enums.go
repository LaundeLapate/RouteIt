package pkg

import "net"

const ServerAddr        string = "127.0.0.1"
const EthernetInterface string = "enx00e04c3607b7"
const LocalInterface    string = "lo"
const WireLessInterface string = "wlp2s0"
const HolePunchPort     uint16 = 8000


const IsPacketForPing  uint8 = 1

var ServerHolePunchAddr net.IPAddr = net.IPAddr{IP: net.ParseIP(ServerAddr)}
