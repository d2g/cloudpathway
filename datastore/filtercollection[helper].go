package datastore

import (
	"github.com/d2g/unqlitego"
	"log"
)

type filtercollectionhelper struct {
	collection              *unqlitego.Database
	domainToCollectionNames *unqlitego.Database
}

var filtercollectionHelperSingleton *filtercollectionhelper = nil

func GetFilterCollectionHelper() (*filtercollectionhelper, error) {
	if filtercollectionHelperSingleton == nil {
		var err error

		filtercollectionHelperSingleton = new(filtercollectionhelper)
		filtercollectionHelperSingleton.collection, err = unqlitego.NewDatabase("FilterCollections.unqlite")
		filtercollectionHelperSingleton.domainToCollectionNames, err = unqlitego.NewDatabase("DomainToFilterCollections.key.unqlite")
		if err != nil {
			return filtercollectionHelperSingleton, err
		}
	}
	return filtercollectionHelperSingleton, nil
}

func (t *filtercollectionhelper) SetFilterCollection(filterCollection FilterCollection) error {
	/*
	 * Get the Existing Collection So we Can Remove any Domain Keys
	 */
	currentCollection, err := t.GetFilterCollection(filterCollection.Name)
	if err != nil {
		return err
	}

	if currentCollection.Name == filterCollection.Name {
		//Do all the deletes in a single Transaction which should make it considerably faster.
		err := t.domainToCollectionNames.Begin()
		if err != nil {
			return err
		}

		// FilterCollection exists
		for _, existingDomain := range currentCollection.Domains {
			//Existing Domains.
			found := false

			for _, domain := range filterCollection.Domains {
				if existingDomain == domain {
					found = true
					break
				}
			}

			if !found {
				//Remove from the Key / Index
				var collectionNames []string
				err = t.domainToCollectionNames.GetObject(existingDomain, &collectionNames)
				if err != nil {

					t.domainToCollectionNames.Rollback()
					return err
				}

				//Delete CollectionName from the array.
				for i, _ := range collectionNames {
					if collectionNames[i] == filterCollection.Name {
						collectionNames[i] = collectionNames[len(collectionNames)-1]
						collectionNames = collectionNames[:len(collectionNames)-1]
						//Save the change to the collection Name
						byteCollectionNames, err := t.domainToCollectionNames.Marshal()(collectionNames)
						if err != nil {
							t.domainToCollectionNames.Rollback()
							return err
						}

						err = t.domainToCollectionNames.Store([]byte(existingDomain), byteCollectionNames)
						if err != nil {
							t.domainToCollectionNames.Rollback()
							return err
						}
						break
					}
				}
			}
		}

		//Commit All our Removals
		err = t.domainToCollectionNames.Commit()
		if err != nil {
			t.domainToCollectionNames.Rollback()
			return err
		}
	}

	//Do all the Adds in a single Transaction which should make it considerably faster.
	err = t.domainToCollectionNames.Begin()
	if err != nil {
		return err
	}

	//Add the Domain Key for the domains that don't already exist.
	for _, domain := range filterCollection.Domains {
		found := false

		//Does it already exists
		if currentCollection.Name == filterCollection.Name {
			for _, existingDomain := range currentCollection.Domains {
				if existingDomain == domain {
					found = true
					break
				}
			}
		}

		if !found {
			//Add the Key / Index
			var collectionNames []string
			err = t.domainToCollectionNames.GetObject(domain, &collectionNames)
			if err != nil {
				return err
			}
			collectionNames = append(collectionNames, filterCollection.Name)

			byteCollectionNames, err := t.domainToCollectionNames.Marshal()(collectionNames)
			if err != nil {
				t.domainToCollectionNames.Rollback()
				return err
			}

			err = t.domainToCollectionNames.Store([]byte(domain), byteCollectionNames)
			if err != nil {
				return err
			}
		}
	}

	//Commit All New Domains
	err = t.domainToCollectionNames.Commit()
	if err != nil {
		t.domainToCollectionNames.Rollback()
		return err
	}

	return t.collection.SetObject(filterCollection.Name, filterCollection)

}

func (this *filtercollectionhelper) GetFilterCollection(name string) (databaseFilterCollection FilterCollection, err error) {
	err = this.collection.GetObject(name, &databaseFilterCollection)
	return
}

func (this *filtercollectionhelper) GetFilterCollections() ([]FilterCollection, error) {
	filterCollections := make([]FilterCollection, 0, 0)

	cursor, err := this.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return filterCollections, err
	}

	err = cursor.First()
	if err != nil {
		//You Get -28 When There are no records.
		if err == unqlitego.UnQLiteError(-28) {
			return filterCollections, nil
		} else {
			return filterCollections, err
		}

	}

	for {
		if !cursor.IsValid() {
			break
		}

		filterCollection := FilterCollection{}
		value, err := cursor.Value()

		if err != nil {

			log.Println("Error: Cursor Get Value Error:" + err.Error())

		} else {

			err := this.collection.Unmarshal()(value, &filterCollection)
			if err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
				} else {
					log.Println("Invalid Filter Collection in Datastore:" + string(key))
					filterCollection.Name = string(key)
					this.SetFilterCollection(filterCollection)
				}
			}

			filterCollections = append(filterCollections, filterCollection)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()

	return filterCollections, err
}

/**
 * Deletes a FilterCollection with the given name.
 */
func (this *filtercollectionhelper) DeleteFilterCollection(name string) error {
	//TODO: Remove the domains index first!
	return this.collection.Delete([]byte(name))
}

func (t *filtercollectionhelper) GetFilterCollectionsWithDomain(domain string) (collections []FilterCollection, err error) {

	var collectionNames []string
	err = t.domainToCollectionNames.GetObject(domain, &collectionNames)
	if err != nil {
		return
	}

	for _, collectionName := range collectionNames {
		var currentCollection FilterCollection
		currentCollection, err = t.GetFilterCollection(collectionName)
		if err != nil {
			return
		}

		if currentCollection.Name == collectionName {
			collections = append(collections, currentCollection)
		}
	}
	return
}
