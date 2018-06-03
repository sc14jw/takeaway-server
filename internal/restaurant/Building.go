package restaurant

// Building represents a singular restaurant building within the system.
type Building struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}
