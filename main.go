package main

import (
	"fmt"
	"log"
	"net/http"
	"takeaway/takeaway-server/internal/vote"

	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
)

func main() {
	voteCtx := &vote.Container{}
	inject.Populate(voteCtx, &vote.MockPollModel{})
	vote.Init(voteCtx)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			vote.GetVote(w, r)
			return
		case http.MethodPut:
			vote.NewVote(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	})

	fmt.Println("Starting server on port 8080. Press ctrl + C to stop it.......")

	log.Fatal(http.ListenAndServe(":8080", r))
}
