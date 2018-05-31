package vote

import (
	"takeaway/takeaway-server/restaurant"
	"testing"
)

var (
	res = &restaurant.Building{
		Name: "r1",
	}
	res2 = &restaurant.Building{
		Name: "r2",
	}
	newRestaurant = &restaurant.Building{
		Name: "New Restaurant",
	}
)

func beforeEach() (p *Poll, empt *Poll) {
	p = &Poll{
		ID:    "poll",
		Votes: createVotes(),
		Options: []*restaurant.Building{
			res,
			res2,
		},
	}
	empt = &Poll{}
	return
}

func createVotes() (v map[string][]string) {
	v = make(map[string][]string)
	v["r1"] = []string{"Jack", "Tom"}
	v["r2"] = []string{"Will", "TJ"}

	return
}

func TestAddVote(t *testing.T) {
	poll, _ := beforeEach()
	poll.AddVote("r1", "test")
	if poll.Votes["r1"][2] != "test" {
		t.Fail()
	}
}

func TestAddVoteNoVotesForOption(t *testing.T) {
	poll, _ := beforeEach()
	poll.AddVote("test", "test")
	if poll.Votes["test"][0] != "test" {
		t.Fail()
	}
}

func TestAddVoteNoVotes(t *testing.T) {
	_, emptyPoll := beforeEach()
	emptyPoll.AddVote("test", "test")
	if emptyPoll.Votes["test"][0] != "test" {
		t.Fail()
	}
}

func TestReAddVote(t *testing.T) {
	poll, _ := beforeEach()
	poll.AddVote("test", "Jack")
	if poll.Votes["test"][0] != "Jack" {
		t.Log("Jack's vote was not assigned to test properly")
		t.Fail()
	} else if len(poll.Votes["r1"]) == 2 {
		t.Log("Jack's vote was not removed from r1 properly")
		t.Fail()
	}
}

func TestClearVotes(t *testing.T) {
	poll, _ := beforeEach()
	poll.ClearVotesFor("Jack")

	if poll.Votes["r1"][0] == "Jack" {
		t.Fail()
	}
}

func TestAddOption(t *testing.T) {
	poll, _ := beforeEach()
	poll.AddOption(newRestaurant)

	if poll.Options[2] != newRestaurant {
		t.Fail()
	}
}

func TestAddOptionNoOptions(t *testing.T) {
	_, empt := beforeEach()
	empt.AddOption(newRestaurant)

	if empt.Options[0] != newRestaurant {
		t.Fail()
	}
}

func TestRemoveOption(t *testing.T) {
	poll, _ := beforeEach()
	poll.RemoveOption(res)

	if poll.Options[0] == res {
		t.Fail()
	} else if poll.Votes[res.Name] != nil {
		t.Fail()
	}
}

func TestRemoveOptionNoOptions(t *testing.T) {
	_, empt := beforeEach()
	empt.RemoveOption(res)

	if empt.Options != nil {
		t.Fail()
	} else if empt.Votes != nil {
		t.Fail()
	}
}

func TestRemoveOptionNotAnOption(t *testing.T) {
	poll, _ := beforeEach()
	poll.RemoveOption(newRestaurant)

	if !restaurantsContains(poll.Options, res) || !restaurantsContains(poll.Options, res2) {
		t.Log("Poll options have been altered after attempting to remove an option that does not exist within the poll")
		t.Fail()
	} else if !stringsContains(poll.Votes[res.Name], "Jack", "Tom") || !stringsContains(poll.Votes[res2.Name], "Will", "TJ") {
		t.Logf("Poll Votes: %v", poll.Votes)
		t.Log("Poll votes have been altered after attempting to remove an option that does not exist within the poll")
		t.Fail()
	}
}

func restaurantsContains(s []*restaurant.Building, elem *restaurant.Building) (found bool) {
	for _, item := range s {
		if item == elem {
			found = true
		}
	}

	return
}

func stringsContains(s []string, elems ...string) (found bool) {
	for _, elem := range elems {
		f := false
		for _, item := range s {
			if item == elem {
				f = true
			}
		}
		if !f {
			return
		}
	}
	found = true
	return
}
