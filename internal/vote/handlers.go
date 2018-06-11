package vote

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"takeaway/takeaway-server/internal/restaurant"
)

const (
	noIDError      = "No ID specified by request from %s."
	unknownIDError = "The given ID %s cannot be found."
)

// locks for polls by ID, ensures that only one goroutine is updating a poll at once.
var pollLocks sync.Map

type vote struct {
	User  string `json:"user"`
	ResID string `json:"restaurant_ID"`
}

// GetPoll provides a http handler for accessing a specified vote.
func GetPoll(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved request")

	// if no ids have been specified within the request, return a bad request status.
	if len(r.URL.Query()["id"]) == 0 {
		log.Println("No ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.URL.Query()["id"][0]
	// if no id is specified as a query parameter, return a bad request status.
	if id == "" {
		log.Println("Empty ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	md := instance.Model
	poll, status, err := md.GetPoll(id)

	// if a poll with the given id could not be found return a status not found response.
	if err != nil {
		log.Printf("Error = %s\n", err.Error())
		log.Printf("Status = %v\n", status)
		if status == NotFound {
			log.Printf("Could not find ID %s, returning not found exception.\n", id)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Unable to find ID due to being unable to connect to the DB.\n")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	data, err := json.Marshal(poll)

	// if poll could not be serialized to JSON, return an internal server error.
	if err != nil {
		log.Printf("The poll %v could not be serialised to JSON.\n", poll)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Returning data: " + string(data[:]))
	_, err = w.Write(data)

	// if poll could not be written to the client, return an internal server error.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

// NewPoll provides a http handler for creating a new vote.
func NewPoll(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if request body could not be parsed, return an internal server error to the client.
	if err != nil {
		log.Println("Could not read body of request")
		http.Error(w, "Could not parse request", http.StatusInternalServerError)
		return
	}

	var data []*restaurant.Building
	err = json.Unmarshal(b, &data)

	// if request could not be properly unmarchelled, return a bad request status to the client.
	if err != nil {
		log.Printf("Could not parse %s into readable object", b)
		http.Error(w, "Could not parse request", http.StatusBadRequest)
		return
	}

	md := instance.Model
	poll, status, err := md.NewPoll(data)

	if err != nil {
		if status == Invalid {
			http.Error(w, "Supplied options invalid", http.StatusBadRequest)
		} else {
			http.Error(w, "Poll could not be created", http.StatusInternalServerError)
		}
		return
	}

	rtnString, err := json.Marshal(poll)
	if err != nil {
		http.Error(w, "Poll could not be created", http.StatusInternalServerError)
		return
	}

	log.Printf("Created poll with id %v\n", poll.ID)
	log.Printf("Returning data: %s\n", rtnString)

	w.WriteHeader(http.StatusCreated)
	w.Write(rtnString)
	return
}

// UpdatePoll allows for a poll within the system to be updated.
func UpdatePoll(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if the body could not be parsed, return an internal server error to the client
	if err != nil {
		log.Println("Could not read body of request")
		http.Error(w, "Could not parse request", http.StatusInternalServerError)
		return
	}

	var data Poll
	err = json.Unmarshal(b, &data)

	// if data cannot be unmarshalled to a Poll object, return a bad request status to the client.
	if err != nil {
		log.Printf("Could not unmarshal passed data into a Poll object data = %s\n", b)
		http.Error(w, "Could not parse request", http.StatusBadRequest)
		return
	}

	md := instance.Model
	status, err := md.UpdatePoll(&data)

	if err != nil {
		if status == NotFound {
			log.Printf("Could not find a poll with specified ID = %s\n", data.ID)
			http.Error(w, "Could not find poll with specified ID", http.StatusBadRequest)
			return
		}

		log.Printf("Could not update poll with id %s due to internal model error %s\n", data.ID, err.Error())
		http.Error(w, "Could not update poll", http.StatusInternalServerError)
		return
	}

	log.Printf("successfully updated poll with id %s\n", data.ID)
	w.WriteHeader(http.StatusAccepted)
}

// DeletePoll allows for a poll to be removed from the system.
func DeletePoll(w http.ResponseWriter, r *http.Request) {
	// if no ids have been specified within the request, return a bad request status.
	if len(r.URL.Query()["id"]) == 0 {
		log.Println("No ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.URL.Query()["id"][0]
	// if no id is specified as a query parameter, return a bad request status.
	if id == "" {
		log.Println("Empty ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	md := instance.Model
	status, err := md.DeletePoll(id)

	if err != nil {
		if status == NotFound {
			// if the given ID could not be found within the datasource return a not found status.
			log.Printf("Could not find the ID %s\n", id)
			http.Error(w, "ID not found", http.StatusNotFound)
		} else {
			// otherwise return an internal server error status.
			http.Error(w, "Could not deal with request", http.StatusInternalServerError)
		}
		return
	}
}

// AddVote provides http handler for adding a vote to a poll.
func AddVote(w http.ResponseWriter, r *http.Request) {
	// if no ids have been specified within the request, return a bad request status.
	if len(r.URL.Query()["id"]) == 0 {
		log.Println("No ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.URL.Query()["id"][0]
	// if no id is specified as a query parameter, return a bad request status.
	if id == "" {
		log.Println("Empty ID specified. Returning bad request status.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	md := instance.Model

	lock := lockPoll(id)
	defer lock.Unlock()

	poll, status, err := md.GetPoll(id)

	// if a poll with the given id could not be found return a status not found response.
	if err != nil {
		log.Printf("Error = %s\n", err.Error())
		log.Printf("Status = %v\n", status)
		if status == NotFound {
			log.Printf("Could not find ID %s, returning not found exception.\n", id)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Printf("Unable to find ID due to being unable to connect to the DB.\n")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// if request body could not be parsed, return an internal server error to the client.
	if err != nil {
		log.Println("Could not read body of request")
		http.Error(w, "Could not parse request", http.StatusInternalServerError)
		return
	}

	var data vote
	err = json.Unmarshal(b, &data)

	// if given body can not be unmarshalled into a vote object, return a bad request status.
	if err != nil {
		log.Printf("Could not unmarshalled %s as a vote\n", data)
		http.Error(w, "Could not parse given vote", http.StatusBadRequest)
		return
	}

	poll.AddVote(data.ResID, data.User)
	status, err = md.UpdatePoll(poll)

	if err != nil {
		if status == NotFound {
			// if the specified poll ID could not be found, return a not found status.
			log.Printf("Could not update poll due to not finding the id %s\n", id)
			http.Error(w, "Could not find poll with specified ID", http.StatusNotFound)
		} else {
			// otherwise return an internal server error status.
			log.Printf("Could not update poll %s due to being unable to connect to the database\n", id)
			http.Error(w, "Could not update poll", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Updated poll %s with a vote for %s for user %s\n", id, data.ResID, data.User)
	w.WriteHeader(http.StatusAccepted)
}

func lockPoll(id string) (lock *sync.Mutex) {
	l, found := pollLocks.Load(id)

	if !found {
		lock = &sync.Mutex{}
		pollLocks.Store(id, lock)
	} else {
		lock = l.(*sync.Mutex)
	}

	lock.Lock()
	return
}
