package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"takeaway/takeaway-server/internal/vote"

	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
)

var (
	useMockData = flag.Bool("useMockData", false, "specify whether or not the server should utilise mock data sources.")
	mongoPort   = flag.Int("mongoPort", 27017, "specify the port the mongo server is currently running on.")
	mongoHost   = flag.String("mongoHost", "localhost", "specify the host the mongo server is currently running on.")
	mongoDB     = flag.String("mongoDB", "takeawayServer", "name of the mongo database where the server should be storing data to.")
)

func main() {
	flag.Parse()

	voteCtx := &vote.Container{}
	if *useMockData {
		log.Println("utilising mock data.")
		inject.Populate(voteCtx, &vote.MockPollModel{})
	} else {
		log.Printf("using mongo instance at %s on port %v\n", *mongoHost, *mongoPort)
		log.Printf("outputting data to %s\n", *mongoDB)
		inject.Populate(voteCtx, &vote.MongoPollModel{
			URL:    *mongoHost + ":" + strconv.Itoa(*mongoPort),
			DBName: *mongoDB,
		})
	}

	vote.Init(voteCtx)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/poll", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			vote.GetPoll(w, r)
		case http.MethodPut:
			vote.NewPoll(w, r)
		case http.MethodPost:
			vote.UpdatePoll(w, r)
		case http.MethodDelete:
			vote.DeletePoll(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	r.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			vote.AddVote(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on port 8080. Press ctrl + C to stop it.......")

	log.Fatal(http.ListenAndServe(":8080", r))
}
