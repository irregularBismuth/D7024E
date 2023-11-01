package src

import (
	"container/list"
	"fmt"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

func (bucket *bucket) ShowContactsInBucket() {
    for e:= bucket.list.Front(); e != nil; e = e.Next(){
        nodeContact := e.Value.(Contact)
        fmt.Println("Node contact in bucket: ", nodeContact.String())
    }
}

func (bucket *bucket) GetLeastRecentlyNode() *Contact {
    bucket_contact := bucket.list.Front().Value.(Contact)
    //head_bucket := bucket.list.Front().Value.(Contact)
    return &bucket_contact
}

func (bucket *bucket) MoveContactTail(target_contact Contact){
    // move contact to tail... 
    for e:= bucket.list.Front(); e != nil; e = e.Next(){
        nodeID := e.Value.(Contact).ID

        if (target_contact).ID.Equals(nodeID){
            bucket.list.MoveToBack(e)
        }
    }
}

func (bucket *bucket) DoesBucketContactExist(target_contact Contact) bool {
    var element *list.Element
    for e:= bucket.list.Front(); e != nil; e = e.Next(){
        nodeID := e.Value.(Contact).ID

        if (target_contact).ID.Equals(nodeID){
            element = e 
        }
    }
    if element != nil {
        return true
    }else{
        return false 
    }
}

func (bucket *bucket) RemoveTargetFromBucket(target_contact Contact){
   
    // iterate from the front to the back and check if target contact matches, if so remove element from bucket.
    var element *list.Element
    for e:= bucket.list.Front(); e != nil; e = e.Next(){
        nodeID := e.Value.(Contact).ID

        if (target_contact).ID.Equals(nodeID){
            element = e 
            //bucket.list.Remove(e)
        }
    }
    if element != nil {
        bucket.list.Remove(element)
    }
}

// GetContactAndCalcDistance returns an array of Contacts where 
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
