package secrets

import "encoding/json"

type IDSet map[string]bool

func (c IDSet) Contains(id string) bool {
	_, ok := c[id]
	return ok
}

func (c IDSet) Add(id string) {
	c[id] = true
}

func (c IDSet) AddAll(ids []string) {
	for _, id := range ids {
		c[id] = true
	}
}

func (c IDSet) Remove(id string) {
	delete(c, id)
}

func (c IDSet) MarshalJSON() ([]byte, error) {
	ids := make([]string, 0, len(c))
	for id := range c {
		ids = append(ids, id)
	}
	return json.Marshal(&ids)
}

func (c *IDSet) UnmarshalJSON(data []byte) error {
	var ids []string

	if err := json.Unmarshal(data, &ids); err != nil {
		return err
	}
	c = &IDSet{}
	c.AddAll(ids)
	return nil
}
