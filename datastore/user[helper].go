package datastore

import (
	"log"

	"github.com/d2g/unqlitego"
)

type userhelper struct {
	collection *unqlitego.Database
}

var userHelperSingleton *userhelper = nil

func GetUserHelper() (*userhelper, error) {
	if userHelperSingleton == nil {
		var err error

		userHelperSingleton = new(userhelper)
		userHelperSingleton.collection, err = unqlitego.NewDatabase("userdata/Users.unqlite")
		if err != nil {
			return userHelperSingleton, err
		}
	}
	return userHelperSingleton, nil
}

func (t *userhelper) GetUser(username string) (databaseUser User, err error) {
	err = t.collection.GetObject(username, &databaseUser)
	return
}

func (t *userhelper) SetUser(user *User) error {
	err := t.collection.SetObject(user.Username, user)
	return err
}

/**
 * Deletes a user with the given username.
 */
func (t *userhelper) DeleteUser(username string) error {
	return t.collection.Delete([]byte(username))
}

/**
 * Return an array of all users.
 */
func (t *userhelper) GetUsers() ([]User, error) {

	users := make([]User, 0, 0)

	cursor, err := t.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return users, err
	}

	err = cursor.First()
	if err != nil {
		if err == unqlitego.UnQLiteError(-28) {
			return users, nil
		} else {
			return users, err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		user := User{}
		value, err := cursor.Value()

		if err != nil {

			log.Println("Error: Cursor Get Value Error:" + err.Error())

		} else {

			err := t.collection.Unmarshal()(value, &user)
			if err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
				} else {
					log.Println("Invalid Lease in Datastore:" + string(key))
					user.Username = string(key)
					t.SetUser(&user)
				}
			}

			users = append(users, user)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()

	return users, err
}

func (t *userhelper) GetUserDisplayName(username string) string {

	databaseUser, err := t.GetUser(username)
	if err != nil {
		return ""
	} else {
		return databaseUser.GetDisplayName()
	}
}
