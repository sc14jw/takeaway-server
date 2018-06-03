package vote

import (
	"sync"
	"takeaway/takeaway-server/internal/restaurant"

	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

var sessionMutex = &sync.Mutex{}

type MongoPollModel struct {
	session *mgo.Session
	DBName  string
	URL     string
}

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
