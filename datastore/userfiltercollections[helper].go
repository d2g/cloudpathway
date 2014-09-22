package datastore

import (
	"github.com/d2g/unqlitego"
	"log"
)

type userfiltercollectionhelper struct {
	collection *unqlitego.Database
}

var userfiltercollectionHelperSingleton *userfiltercollectionhelper = nil

func GetUserFilterCollectionsHelper() (*userfiltercollectionhelper, error) {
	if userfiltercollectionHelperSingleton == nil {
		var err error

		userfiltercollectionHelperSingleton = new(userfiltercollectionhelper)
		userfiltercollectionHelperSingleton.collection, err = unqlitego.NewDatabase("UsersFilterCollections.unqlite")
		if err != nil {
			return userfiltercollectionHelperSingleton, err
		}
	}
	return userfiltercollectionHelperSingleton, nil
}

func (this *userfiltercollectionhelper) SetUserFilterCollections(userFilterCollections UserFilterCollections) error {
	return this.collection.SetObject(userFilterCollections.Username, userFilterCollections)
}

func (this *userfiltercollectionhelper) GetUserFilterCollections(username string) (databaseUserFilterCollections UserFilterCollections, err error) {
	err = this.collection.GetObject(username, &databaseUserFilterCollections)

	// Set the username on the object, if we loaded one successfully, this will make no difference. If
	// the record didn't exist and we have a new instance, this will set the key.
	databaseUserFilterCollections.Username = username

	return databaseUserFilterCollections, err
}

func (t *userfiltercollectionhelper) GetUsersFilterCollections() ([]UserFilterCollections, error) {
	usersFilterCollections := make([]UserFilterCollections, 0, 0)

	cursor, err := t.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return usersFilterCollections, err
	}

	err = cursor.First()
	if err != nil {
		if err == unqlitego.UnQLiteError(-28) {
			return usersFilterCollections, nil
		} else {
			return usersFilterCollections, err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		userFilterCollections := UserFilterCollections{}
		value, err := cursor.Value()

		if err != nil {

			log.Println("Error: Cursor Get Value Error:" + err.Error())

		} else {

			err := t.collection.Unmarshal()(value, &userFilterCollections)
			if err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
				} else {
					log.Println("Invalid Lease in Datastore:" + string(key))
					userFilterCollections.Username = string(key)
					t.SetUserFilterCollections(userFilterCollections)
				}
			}

			usersFilterCollections = append(usersFilterCollections, userFilterCollections)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()

	return usersFilterCollections, err
}
