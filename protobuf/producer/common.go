package producer

import "fmt"

// GetNodeKey returns the key used by the Node's Topic based on the NodeCriteria
func GetNodeKey(c *NodeCriteria) string {
	var key string
	if c.ForeignSource != "" && c.ForeignId != "" {
		key = fmt.Sprintf("%s:%s", c.ForeignSource, c.ForeignId)
	} else {
		key = fmt.Sprintf("%d", c.Id)
	}
	return key
}
