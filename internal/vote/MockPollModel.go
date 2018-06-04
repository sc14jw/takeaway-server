package vote

import (
	"fmt"

	"takeaway/takeaway-server/internal/restaurant"
)

var r1 = &restaurant.Building{
	ID:      "r1",
	Name:    "Restaurant 1",
	Address: "Address 1",
}

var r2 = &restaurant.Building{
	ID:      "r2",
	Name:    "Restaurant 2",
	Address: "Address 2",
}

// MockPollModel provides a mock implementation of the PollModel interface.
type MockPollModel struct {
	p *Poll
}

// GetPoll returns a saved poll with the id "test", or nil with an error should the ID "unknown" be passed to the method.
func (pm *MockPollModel) GetPoll(id string) (poll *Poll, status Status, err error) {
	if id == "unknown" {
		err = fmt.Errorf("the ID %s is not a valid poll ID", id)
		status = NotFound
		return
	}

	if pm.p == nil {
		v := make(map[string][]string)
		v["r1"] = []string{"Jack", "Tom"}
		v["r2"] = []string{"Will", "TJ"}

		opts := []*restaurant.Building{r1, r2}

		pm.p = &Poll{
			ID:      "test",
			Votes:   v,
			Options: opts,
		}
	}
	poll = pm.p
	return
}

// NewPoll creates a new poll returning the created poll. This poll is used as the saved poll for the mock. An error will be returned from this method should the first option's name passed be "unknown", returning nil
// for the returned poll.
func (pm *MockPollModel) NewPoll(options []*restaurant.Building) (poll *Poll, status Status, err error) {
	if options[0].Name == "unknown" {
		err = fmt.Errorf("the specified option %s could not be found", options[0])
		status = NotFound
		return
	}

	poll = &Poll{
		ID:      "new poll",
		Options: options,
	}

	pm.p = poll

	return
}

// Close has been added to ensure the mock meets the PollModel interface, it does not need to actually complete anything.
func (pm *MockPollModel) Close() (err error) {
	return
}