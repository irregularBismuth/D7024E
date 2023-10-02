package src

import (
	"net"
    "fmt"
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
    fmt.Printf("New node was created with: \n Address: %s\n Contact: %s\n ID: %s\n",address.String(), new_contact.String(), id_node.String()) 
    return Kademlia{
        node_contact: routing_table,
        data: make(map[string]string),
    }
}
// Node lookup algorithm 
func (kademlia *Kademlia) LookupContact(node_network *Network, target *Contact) []Contact {
	// *Node position = shortest unique prefix in tree
    // 1. Pick nodes for the closest non-empty k-bucket
    // 2. Send parallel async "FIND_NODE" RPCs

    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
    // If bucket is not full, the new contact is added at the end. 

    shortlist_contacts := kademlia.node_contact.FindClosestContacts(target.ID, 3)
    //contacted_nodes := ContactCandidates{}

    for i := 0; i < len(shortlist_contacts); i++ {
        contact := shortlist_contacts[i]
        //println(contact.String())
        
        // call RPC for FIND_NODE here
        // FIND NODE = k-closest = []Contacts
        target_addr,_ := net.ResolveUDPAddr("udp", contact.Address)
        response := node_network.FetchRPCResponse(FindNode,"lookup_rpc_id",&contact, target_addr)
        fmt.Println("FIND NODE response: ",response.Contacts)
        
    }
    return shortlist_contacts

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
