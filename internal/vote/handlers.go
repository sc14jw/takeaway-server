package vote

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"takeaway/takeaway-server/internal/restaurant"
)

const (
	noIDError      = "No ID specified by request from %s."
	unknownIDError = "The given ID %s cannot be found."
)

// GetVote provides a http handler for accessing a specified vote.
func GetVote(w http.ResponseWriter, r *http.Request) {
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

// NewVote provides a http handler for creating a new vote.
func NewVote(w http.ResponseWriter, r *http.Request) {
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