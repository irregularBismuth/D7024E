package src

import (
    "net"
)


type Network struct {
    
}

func Listen(ip string, port int) {
	// TODO
    var udp = net.UDPAddr{
        IP: net.IP(ip),
        Port: port,
    }
    
    net.ListenUDP(ip, &udp)
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
