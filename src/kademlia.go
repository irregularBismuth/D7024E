package src

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
)

// Kademlia nodes store contact information about each other <IP, UDP port, Node ID>
type Kademlia struct {
    node_contact RoutingTable
    data map[string][]byte
}

func InitNode(address net.Addr) Kademlia {
    var id_node *KademliaID = NewRandomKademliaID()
    var new_contact Contact = NewContact(id_node, address.String())
    var routing_table RoutingTable = *NewRoutingTable(new_contact)
    fmt.Println("New node was created with Address: ",address.String()) 
    return Kademlia{
        node_contact: routing_table,
        data: make(map[string][]byte),
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

func (kademlia *Kademlia) LookupData(hash string) ([]byte, bool){
    // Take the sha1 encryption and check if it exists as a key
    //fmt.Println("1. Hash to look up: " + hash)
    original, exists := kademlia.data[hash]
    //fmt.Printf("Original = %x", original)
    fmt.Println("")
    if exists{
        fmt.Printf("The data you want exists: %s", original)
        fmt.Println("")
    } else{
        fmt.Println("The data does not exist")
    }
    return original, exists
}

func (kademlia *Kademlia) Store(data []byte) {
	// Encrypt the hash for our value
    hash := Hash(data)

    // Save the key value pair to current node
    kademlia.data[hash] = data
    //fmt.Println("Storing key value pair, DONE")
    fmt.Println("Stored the hash: " + hash)
    //fmt.Println("Key value is: " + kademlia.data[hash])
    return
}

func Hash(data []byte) (string) {
    // Create the hash value
    hasher := sha1.Sum(data)

    // Convert the hash to hexadecmial string
    hash := hex.EncodeToString(hasher[0:IDLength])
    fmt.Println("Hashing DONE")
    return hash
}
