package src

import ("net")

// Kademlia nodes store contact information about each other <IP, UDP port, Node ID>
type Kademlia struct {
    node_contact Contact
    contact_table RoutingTable
    data map[string]string
}

func InitNode(ip net.IP) Kademlia {
    
    var id_node *KademliaID = NewRandomKademliaID()
    var new_contact Contact = NewContact(id_node, ip.String())
    var routing_table RoutingTable = *NewRoutingTable(new_contact)
    
    return Kademlia{
        node_contact: new_contact,
        contact_table: routing_table,
        data: make(map[string]string),
    }
    
    
}




func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
