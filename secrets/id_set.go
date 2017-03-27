package secrets

import (
	"bytes"
	"encoding/json"
)

// IDSet is a helper to handie a set of generic ids
type IDSet map[string]bool

// Contains checks if a given id is part of the set
func (c IDSet) Contains(id string) bool {
	_, ok := c[id]
	return ok
}

// Add a given id to the set
func (c IDSet) Add(id string) {
	c[id] = true
}

// AddAll adds a slice of ids to the set
func (c IDSet) AddAll(ids []string) {
	for _, id := range ids {
		c[id] = true
	}
}

// Remove a given id from the set
func (c IDSet) Remove(id string) {
	delete(c, id)
}

// Equals compares this id set to another for equality
func (c IDSet) Equals(other IDSet) bool {
	for id := range c {
		if !other.Contains(id) {
			return false
		}
	}
	for id := range other {
		if !c.Contains(id) {
			return false
		}
	}
	return true
}

func (c IDSet) String() string {
	buf := bytes.NewBufferString("[")
	delim := ""
	for id := range c {
		buf.WriteString(delim)
		buf.WriteString(id)
		delim = " "
	}
	buf.WriteString("]")
	return buf.String()
}

// MarshalJSON creates a json array of the set
func (c IDSet) MarshalJSON() ([]byte, error) {
	ids := make([]string, 0, len(c))
	for id := range c {
		ids = append(ids, id)
	}
	return json.Marshal(&ids)
}

// UnmarshalJSON reads the set from a json array
func (c *IDSet) UnmarshalJSON(data []byte) error {
	var ids []string

	if err := json.Unmarshal(data, &ids); err != nil {
		return err
	}
	idSet := IDSet{}
	idSet.AddAll(ids)
	*c = idSet
	return nil
}
