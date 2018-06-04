package vote

import (
	"sync"
	"takeaway/takeaway-server/internal/restaurant"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

var sessionMutex = &sync.Mutex{}

// MongoPollModel provides a mongo based implementation to the PollModel interface.
type MongoPollModel struct {
	session *mgo.Session
	DBName  string
	URL     string
}

// GetPoll gets a poll from the mongo database with the specified id, returning the found poll as a Poll object, a status
// and an error should any issues occur while trying to return the specified poll.
func (pm *MongoPollModel) GetPoll(id string) (poll *Poll, status Status, err error) {
	err = pm.openSessionIfRequired()
	if err != nil {
		status = NoConnection
		return
	}

	p := Poll{}

	c := pm.session.DB(pm.DBName).C("polls")
	err = c.Find(bson.M{"id": id}).One(&p)

	if err != nil {
		status = NotFound
		return
	}

	poll = &p

	return
}

// NewPoll creates a new poll within the mongo database, returning the created Poll object with a status and any errors
// that occur while attempting to create the poll.
func (pm *MongoPollModel) NewPoll(options []*restaurant.Building) (poll *Poll, status Status, err error) {
	err = pm.openSessionIfRequired()
	if err != nil {
		status = NoConnection
		return
	}

	data := Poll{
		ID:      bson.NewObjectId().Hex(),
		Options: options,
	}

	poll = &data

	c := pm.session.DB(pm.DBName).C("polls")
	err = c.Insert(data)

	if err != nil {
		status = Invalid
		return
	}

	return
}

// Close allows the model to be closed properly, ensuring any mongo sessions are properly closed.
func (pm *MongoPollModel) Close() (err error) {
	err = pm.Close()
	return
}

func (pm *MongoPollModel) openSessionIfRequired() (err error) {
	sessionMutex.Lock()
	if pm.session == nil {
		pm.session, err = mgo.Dial(pm.URL)
	}
	sessionMutex.Unlock()
	return
}
