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

    fmt.Printf("New node was created with: \n Address: %s\n Contact: %s\n ID: %s\n",address.String(), new_contact.String(), id_node.String()) 
    return Kademlia{
        node_contact: routing_table,
        data: make(map[string][]byte),
    }
}
// Node lookup algorithm 
func (kademlia *Kademlia) LookupContact(node_network *Network, target *Contact) {
	// *Node position = shortest unique prefix in tree
    // 1. Pick nodes for the closest non-empty k-bucket
    // 2. Send parallel async "FIND_NODE" RPCs
   
    // STEP 1 - Get the k-closest from known nodes 
    alpha := 3
    contacted_nodes := ContactCandidates{}
    shortlist_contacts := kademlia.node_contact.FindClosestContacts(target.ID, alpha)
   
    // STEP 2 - Update the shortlist contacts based on shortest distance from the k-closest response contacts 
    for i := 0; i < len(shortlist_contacts); i++ {
        contact := shortlist_contacts[i] //known contact 
        target_addr,_ := net.ResolveUDPAddr("udp", contact.Address)
       
        // STEP 3 - 'FIND_NODE' RPC call will return k-closest nodes to the target node 
        response, request_error := node_network.FetchRPCResponse(FindNode,"lookup_rpc_id",target,target_addr)
        
        if request_error != nil {
            fmt.Println("Request error: ", request_error.Error())
        }else{
            // STEP 4 - Update the shortlist (k-closest contacts) based on the received response... 
            fmt.Println("FIND NODE response: ",response.Contacts)
            shortlist_contacts = response.Contacts // update shortlist contacts initial value from the first response 
        } 

        // STEP 5 - Contact the received k-contacts from the response
        for i := 0; i < len(shortlist_contacts); i++{
            k_contact := shortlist_contacts[i]
            k_addr, _ := net.ResolveUDPAddr("udp", k_contact.Address)
            k_response, k_request_error := node_network.FetchRPCResponse(FindNode,"lookup_rpc_id", target, k_addr)
                
            if k_request_error != nil {
                fmt.Println("Request error: ", k_request_error.Error())

            }else{
                fmt.Println("FIND NODE response: ",k_response.Contacts)
                
                // STEP 6 - update shortlist if distance is less than prior + append contacted nodes in the ContactCandidates array
                for _, c_response := range k_response.Contacts {
                    // Iterate the current shortlist and check if distance is less and update accordingly, need to also check if it has been contacted 
                    for i_short, shortlist_contact := range shortlist_contacts{  
                        if(c_response.distance.Less(shortlist_contact.ID)){
                            if (len(shortlist_contacts) < alpha){
                                shortlist_contacts = append(shortlist_contacts, c_response) 
                            }else{
                                shortlist_contacts[i_short] = c_response 
                            }
                        }else{
                            continue
                        }
                    }
                }


            }
                
        }
        
    }

}

func (kademlia *Kademlia) LookupData(network *Network,hash string) ([]byte, bool){
    // Take the sha1 encryption and check if it exists as a key
    fmt.Println("1. Hash to look up: " + hash)
    original, exists := kademlia.data[hash]
    fmt.Printf("Original = %x", original)
    fmt.Println("")
    var a []Contact 
    if exists{
        fmt.Printf("2. The data you want exists: %s", original)
        fmt.Println("")
    } else{
        fmt.Println("2. Does not exist")
        fmt.Println("The data can't be found on self, searching through K closest Nodes")
        boot_addr, _ := GetBootnodeAddr()
        bootnode_contact := network.FetchRPCResponse(Ping,"boot_node_contact_id",&network.node.node_contact.me,boot_addr,"")
        
        a = kaddemlia.node.node_contact.AddContact(&bootnode_contact.Contact)
        for i:=0; i<len(a); i++{
            target_node := a[i]
            target_addr, _ := net.ResolveUDPAddr("udp",target_node.Address)
            value_response := network.FetchRPCResponse(FindValue,"",&network.node.node_contact.me,target_addr,hash)
            kademlia.Store(value_response.Value)
            fmt.Printf("Object found and downloaded %s",string(value_response.Value))
        }
        
        return original,exists

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
