package src

// Kademlia nodes store contact information about each other <IP, UDP port, Node ID>
type Kademlia struct {
    node_id string
    contact *RoutingTable
    data map[string]string
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
