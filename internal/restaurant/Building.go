package restaurant

// Building represents a singular restaurant building within the system.
type Building struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
