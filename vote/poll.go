package vote

import "takeaway/takeaway-server/restaurant"

// Poll represents a singular vote within the system.
type Poll struct {
	ID      string                 `json:"id"`
	Votes   map[string][]string    `json:"votes"`
	Options []*restaurant.Building `json:"options"`
}
