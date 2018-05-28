package vote

import (
	"fmt"
	"takeaway/takeaway-server/restaurant"
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
func (pm *MockPollModel) GetPoll(id string) (poll *Poll, err error) {
	if id == "unknown" {
		err = fmt.Errorf("the ID %s is not a valid poll ID", id)
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
