package vote

import (
	"encoding/json"
	"log"
	"net/http"
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
	poll, err := md.GetPoll(id)

	// if a poll with the given id could not be found return a status not found response.
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(poll)

	// if poll could not be serialized to JSON, return an internal server error.
	if err != nil {
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
