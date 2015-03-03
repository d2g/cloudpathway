package datastore

import "github.com/d2g/unqlitego"

type domaincollectionhelper struct {
	collections             map[string]*unqlitego.Database
	domainToCollectionNames *unqlitego.Database
}

var domaincollectionhelperSingleton *domaincollectionhelper = nil

func GetDomainCollectionHelper() (*domaincollectionhelper, error) {
	if domaincollectionhelperSingleton == nil {
		var err error

		domaincollectionhelperSingleton = new(domaincollectionhelper)
		domaincollectionhelperSingleton.collections = make(map[string]*unqlitego.Database)
		domaincollectionhelperSingleton.domainToCollectionNames, err = unqlitego.NewDatabase("userdata/DomainToDomainCollections.key.unqlite")
		if err != nil {
			return domaincollectionhelperSingleton, err
		}
	}
	return domaincollectionhelperSingleton, nil
}

func (t *domaincollectionhelper) GetDomainCollections() ([]DomainCollection, error) {
	return []DomainCollection{}, nil
}

func (t *domaincollectionhelper) GetDomainCollection(name string) (DomainCollection, error) {
	return DomainCollection{}, nil
}

func (t *domaincollectionhelper) SetDomainCollection(dc DomainCollection) error {
	return nil
}

/**
 * Deletes a FilterCollection with the given name.
 */
func (t *domaincollectionhelper) DeleteDomainCollection(name string) error {
	//TODO: Remove the domains index first!
	return nil
}

func (t *domaincollectionhelper) GetDomainCollectionsWith(domain string) (collections []DomainCollection, err error) {

	err = t.domainToCollectionNames.GetObject(domain, &collections)
	return
}
