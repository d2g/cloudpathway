package datastore

import ()

type UserFilterCollections struct {
	Username    string
	Collections []string //Array Of Collection Names...
}

/**
 * Returns the number of collections.
 */
func (this *UserFilterCollections) NumberOfCollections() int {
	if this.Collections == nil {
		return 0
	} else {
		return len(this.Collections)
	}
}

func (t *UserFilterCollections) ContainsCollection(searchCollectionName string) bool {
	for _, collectionName := range t.Collections {
		if collectionName == searchCollectionName {
			return true
		}
	}
	return false
}