package vote

var instance *Container

// PollModel defines a contract for how the system should interact with the database for accessing poll information.
type PollModel interface {
	// GetPoll allows for a singular poll to be accessed, using its ID. Should any issue occur while attempting to access the poll specified by the ID, an error will be returned. Should a poll be
	// located using the specified ID, the poll will be returned as a pointer to a Poll object.
	GetPoll(id string) (*Poll, error)
}

// Container provides access to injected implementation of PollModel for the application.
type Container struct {
	Model PollModel `inject:""`
}

// Init allows the vote package to be initialised with the Container c.
func Init(c *Container) {
	instance = c
}
