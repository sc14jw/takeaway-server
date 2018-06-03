package vote

import "takeaway/takeaway-server/internal/restaurant"

// Poll represents a singular vote within the system.
type Poll struct {
	ID      string                 `json:"id" bson:"id"`
	Votes   map[string][]string    `json:"votes" bson:"votes"`
	Options []*restaurant.Building `json:"options" bson:"options"`
}

// AddOption allows for a restaurant to be added to the poll object.
func (p *Poll) AddOption(opt *restaurant.Building) {
	if p.Options == nil {
		p.Options = make([]*restaurant.Building, 0)
	}
	p.Options = append(p.Options, opt)
}

// AddVote allows a singular vote to be added for a given option for a specified user. Should this user have voted in for any other option previously, their previous vote will be removed.
func (p *Poll) AddVote(opt string, name string) {
	if p.Votes == nil {
		p.Votes = make(map[string][]string)
	}

	if p.Votes[opt] == nil {
		p.Votes[opt] = make([]string, 0)
	}

	p.ClearVotesFor(name)
	p.Votes[opt] = append(p.Votes[opt], name)
}

// ClearVotesFor allows for the votes for a given user to be removed from the poll.
func (p *Poll) ClearVotesFor(user string) {
	for k := range p.Votes {
		for i, n := range p.Votes[k] {
			if n == user {
				p.Votes[k] = append(p.Votes[k][:i], p.Votes[k][i+1:]...)
			}
		}
	}
}

// RemoveOption allows for a given restaurant to be removed as an option within the poll. This does mean any votes currently cast for the given restaurant will be lost.
func (p *Poll) RemoveOption(restaurant *restaurant.Building) {
	for i, elem := range p.Options {
		if elem == restaurant {
			p.Options = append(p.Options[:i], p.Options[i+1:]...)
		}
	}

	if p.Votes != nil {
		p.Votes[restaurant.Name] = nil
	}
}
