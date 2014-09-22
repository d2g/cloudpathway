package datastore

import (
	"html/template"
	"strings"
)

type FilterCollection struct {
	Name    string   //Collection Name
	Domains []string // Usually domains but should also be ok for IPs
}

/**
 * Returns the number of domains on this collection.
 */
func (this *FilterCollection) GetNumberOfDomains() int {
	return len(this.Domains)
}

/**
 * Returns a HTML safe escaped version of the name of this collection.
 */
func (this *FilterCollection) EscapedName() string {
	result := strings.Replace(this.Name, " ", "+", -1)
	return template.JSEscapeString(result)
}
