package datastore

import (
//"html/template"
//"log"
//"strings"
)

type DomainCollection struct {
	Name string

	Domains DomainNames
}

// Paged Domain Names with Capacity
type DomainNames struct {
	values []string

	length int
}

func (t DomainNames) Length() (int, error) {
	return 1, nil
}

/**
 * Returns the number of domains on this collection.
 */
func (t *DomainCollection) GetNumberOfDomains() int {

	return 1
}

/**
 * Returns a HTML safe escaped version of the name of this collection.
 */
func (t *DomainCollection) EscapedName() string {
	//result := strings.Replace(string(*t), " ", "+", -1)
	//return template.JSEscapeString(result)
	return ""
}
