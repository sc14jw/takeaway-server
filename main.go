package main

import (
	"fmt"
	"log"
	"net/http"
	"takeaway/takeaway-server/vote"

	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
)

func main() {
	voteCtx := &vote.Container{}
	inject.Populate(voteCtx, &vote.MockPollModel{})
	vote.Init(voteCtx)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/vote", vote.GetVote)

	fmt.Println("Starting server on port 8080. Press ctrl + C to stop it.......")

	log.Fatal(http.ListenAndServe(":8080", r))
}
