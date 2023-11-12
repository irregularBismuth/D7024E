package src

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	//"sync"
)

// Kademlia nodes store contact information about each other <IP, UDP port, Node ID>
type Kademlia struct {
    NetworkInterface NetworkInterface
    node_contact RoutingTable
    data map[string][]byte
}

func InitNode(address *net.UDPAddr) Kademlia {
    //address, _ := GetLocalAddr()
    var id_node *KademliaID = NewRandomKademliaID()
    var new_contact Contact = NewContact(id_node, address.String())
    var routing_table RoutingTable = *NewRoutingTable(new_contact)
    //var new_network Network = InitNetwork()

    fmt.Printf("New node was created with: \n Address: %s\n Contact: %s\n ID: %s\n",address.String(), new_contact.String(), id_node.String()) 
    kademlia := Kademlia{
        //NetworkInterface: network,
        node_contact: routing_table,
        data: make(map[string][]byte),
    }
    return kademlia
}

func (kademlia *Kademlia) SetNetworkInterface(network NetworkInterface){
    kademlia.NetworkInterface = network
}

func (kademlia *Kademlia) ShowNodeBucketStatus(){
    //current_buckets := network.node.node_contact.buckets
    //known_buckets := network.node.node_contact.buckets
    known_buckets := kademlia.node_contact.buckets
    for i:=0; i < len(known_buckets); i++{
        bucket := known_buckets[i]
        if bucket.Len() > 0 {
            bucket.ShowContactsInBucket()
        }
    }
}


// This function update appropriate k-bucket for the sender's node ID.
// The argument takes the target contacts bucket received from requests or response
func (kademlia *Kademlia) UpdateHandleBuckets(target_contact Contact){
   
    // Fetch the correct bucket location based on bucket index
    bucket_index := kademlia.node_contact.getBucketIndex(target_contact.ID)
    target_bucket := kademlia.node_contact.buckets[bucket_index]
    
    // if bucket is not full = add the node to the bucket 
    if target_bucket.Len() < GetMaximumBucketSize() && target_bucket.DoesBucketContactExist(target_contact) {
        fmt.Printf("Bucket contact already exist adding the contact to tail: %s\n", target_contact.Address)
        kademlia.node_contact.AddContact(target_contact)
    
    }else if target_bucket.Len() == GetMaximumBucketSize() {
        // If bucket is full - ping the k-bucket's least-recently seen node
        // least-recently node at the head & most-recently node at the tail 
       
        least_recently_node := target_bucket.GetLeastRecentlyNode() // contact at the tail
        least_recently_addr, _ := net.ResolveUDPAddr("udp", least_recently_node.Address)
        
        my_contact := kademlia.node_contact.me
        fmt.Printf("Bucket was full trying to ping recently-seen node: %s\n",least_recently_node.Address)
        rpc_ping, request_err := kademlia.NetworkInterface.FetchRPCResponse(Ping, "bucket_full_ping_id", &my_contact, least_recently_addr)
       
        // Might not need this rpc_ping handling since every rpc will still run the 'UpdateHandleBuckets method'

        // Todo: Maybe add a new branch for checking the ping for the k-bucket's least_recently_node

        if request_err != nil || rpc_ping.Error != nil {
            //failed to response - removed from the k-bucket and new sender inserted at the tail
            fmt.Println("Request was unsuccessful, removing least-recently seen node from bucket: ", least_recently_node.Address)
            kademlia.node_contact.RemoveTargetContact(*least_recently_node)
            
        }else if rpc_ping.Error == nil{
            //successful response -  contact is moved to the tail of the list 
            fmt.Printf("Ping was successful adding/moving contact: %s to tail\n", rpc_ping.Contact.Address)
            kademlia.node_contact.AddContact(rpc_ping.Contact)
        }

    }else {
        kademlia.node_contact.AddContact(target_contact)        
    }

    kademlia.ShowNodeBucketStatus()
}

func (kademlia *Kademlia) AsynchronousLookupContact(target_contact Contact){
    alpha := 3
    contacted_nodes := ContactCandidates{contacts: []Contact{}}
    result_shortlist := ContactCandidates{kademlia.node_contact.FindClosestContacts(target_contact.ID, alpha)}
    response_channel := make(chan PayloadData, alpha) // Create a temporary response channel of size alpha for contact and contacts handling
  
    for len(contacted_nodes.contacts) < alpha && result_shortlist.Len() > 0{
        
        //var mutex sync.Mutex
        //var wait_group sync.WaitGroup

        for i := 0; i < result_shortlist.Len() && len(contacted_nodes.contacts) < alpha; i++ {
            
            contact := result_shortlist.contacts[i] //known contact 
            target_addr,_ := net.ResolveUDPAddr("udp", contact.Address) 
            go kademlia.NetworkInterface.AsynchronousFindNode(target_contact, target_addr, response_channel)
        
        }
        // TO FIX:
        // 1. Make so it updates already contacted
        // 2. Update shortlist      
        
        k_response := <- response_channel
        k_closest := k_response.Contacts
        k_contact := k_response.Contact
        if len(k_closest) > 0 {
            contacted_nodes.contacts = append(contacted_nodes.contacts, k_contact) // Add contacted node       
        }
         
    }
}


// Node lookup algorithm 
func (kademlia *Kademlia) LookupContact(target *Contact) {
   
    // STEP 1 - Get the k-closest from known nodes 
    alpha := 3
    contacted_nodes := ContactCandidates{}
    shortlist_contacts := kademlia.node_contact.FindClosestContacts(target.ID, alpha)
   
    // STEP 2 - Update the shortlist contacts based on shortest distance from the k-closest response contacts 
    for i := 0; i < len(shortlist_contacts); i++ {
        contact := shortlist_contacts[i] //known contact 
        target_addr,_ := net.ResolveUDPAddr("udp", contact.Address)
       
        // STEP 3 - 'FIND_NODE' RPC call will return k-closest nodes to the target node 
        response, request_error := kademlia.NetworkInterface.FetchRPCResponse(FindNode,"lookup_rpc_id",target,target_addr)
         
        if request_error != nil {
            fmt.Println("Request error: ", request_error.Error())
        }else{
            // STEP 4 - Update the shortlist (k-closest contacts) based on the received response... 
            contacted_nodes.contacts = append(contacted_nodes.contacts, response.Contact)
            fmt.Println("FIND NODE response: ",response.Contacts)

            // TODO write a method for checking a list of contacts and update based on shortest distance
            
            //shortlist_contacts = response.Contacts // update shortlist contacts initial value from the first response 
        } 

        // STEP 5 - Contact the received k-contacts from the response
        for i := 0; i < len(shortlist_contacts); i++{
            k_contact := shortlist_contacts[i]
            //k_contact := response.Contacts[i]
            
            k_addr, _ := net.ResolveUDPAddr("udp", k_contact.Address)
            k_response, k_request_error := kademlia.NetworkInterface.FetchRPCResponse(FindNode,"lookup_rpc_id", target, k_addr)
                
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

/*
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
        bootnode_contact := Network.FetchRPCResponse(Ping,"boot_node_contact_id",&network.node.node_contact.me,boot_addr,"")
        
        a = kademlia.node.node_contact.AddContact(&bootnode_contact.Contact)
        for i:=0; i<len(a); i++{
            target_node := a[i]
            target_addr, _ := net.ResolveUDPAddr("udp",target_node.Address)
            value_response := Network.FetchRPCResponse(FindValue,"",&network.node.node_contact.me,target_addr,hash)
            kademlia.Store(value_response.Value)
            fmt.Printf("Object found and downloaded %s",string(value_response.Value))
        }
        
        return original,exists

    }

    return original, exists
}
*/

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
