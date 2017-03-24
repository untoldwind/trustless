package secrets

import "encoding/json"

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
	c = &IDSet{}
	c.AddAll(ids)
	return nil
}
