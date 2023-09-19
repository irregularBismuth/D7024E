package src

import (
	"net"
)

// Kademlia nodes store contact information about each other <IP, UDP port, Node ID>
type Kademlia struct {
    node_contact RoutingTable
    data map[string]string
}

func InitNode(address net.Addr) Kademlia {
    var id_node *KademliaID = NewRandomKademliaID()
    var new_contact Contact = NewContact(id_node, address.String())
    var routing_table RoutingTable = *NewRoutingTable(new_contact)
    //fmt.Println("New node was created with ID: ",id_node.String()) 
    return Kademlia{
        node_contact: routing_table,
        data: make(map[string]string),
    }
}
// Node lookup algorithm 
func (kademlia *Kademlia) LookupContact(target *Contact) {
	// *Node position = shortest unique prefix in tree
    // 1. Pick nodes for the closest non-empty k-bucket
    // 2. Send parallel async "FIND_NODE" RPCs

    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
    // If bucket is not full, the new contact is added at the end. 

    shortlist_closest_contacts := kademlia.node_contact.FindClosestContacts(target.ID, 3)
    //contact_list := ContactCandidates{}
    //kademlia.contact_table.buckets

    for i := 0; i < len(shortlist_closest_contacts); i++ {
        contact := shortlist_closest_contacts[i]
        println(contact.String())
        // call RPC for FIND_NODE here
        
        
    }

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
