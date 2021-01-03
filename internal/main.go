/*
Copyright 2020 Joseph Davis

For license infomation (MIT), see the containing github project.

As can be seen below, this version uses Gorilla Mux. It can be installed with:
go get github.com/gorilla/mux
And PostgreSQL with
go get -u github.com/lib/pq

I used https://dev.to/moficodes/build-your-first-rest-api-with-go-2gcj
as a starting point.
And added in https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
for database and testing.

Including running a basic postgres instance with 
  sudo docker run -it -p 5432:5432 -d postgres
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)


// WelcomeGet - A simple landing page response to give a few hints about usage
func WelcomeGet(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// TODO: include more helpful information about the API
    w.Write([]byte(`{"message": "Welcome to the Worst3D Printer Slicing Service!"}`))
}

// HealthCheckGet - simple health check endpoint, just returns 200
func HealthCheckGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Running"}`))
}

// Create a new user, if the request is authenticated
func UserCreate(w http.ResponseWriter, r *http.Request) {

}


/* Basic functions */
func post(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"message": "post called"}`))
}

func put(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(`{"message": "put called"}`))
}

func delete(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "delete called"}`))
}

func params(w http.ResponseWriter, r *http.Request) {
    pathParams := mux.Vars(r)
    w.Header().Set("Content-Type", "application/json")

    userID := -1
    var err error
    if val, ok := pathParams["userID"]; ok {
        userID, err = strconv.Atoi(val)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"message": "need a number"}`))
            return
        }
    }

    commentID := -1
    if val, ok := pathParams["commentID"]; ok {
        commentID, err = strconv.Atoi(val)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"message": "need a number"}`))
            return
        }
    }

    query := r.URL.Query()
    location := query.Get("location")

    w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
}

func main() {
	/* quick site map
	/worst3d/v3
	  /user
	  /job
	  /project
	  /healthcheck
	*/

	//r := mux.NewRouter()
	worstApp := Worst3DPSApp{}
	worstApp.Initialize(
		os.Getenv("WORST_DB_USERNAME"),
		os.Getenv("WORST_DB_PASSWORD"),
		os.Getenv("WORST_DB_NAME") )

	//r := worstApp.Router

	api := worstApp.Router.PathPrefix("/worst3d/v3").Subrouter()

	api.HandleFunc("", WelcomeGet).Methods(http.MethodGet)
	api.HandleFunc("/healthcheck", HealthCheckGet).Methods(http.MethodGet)
	
	api.HandleFunc("/user", params).Methods(http.MethodPost)  // user creation
    api.HandleFunc("/user/{userID}", params).Methods(http.MethodGet)

    api.HandleFunc("", post).Methods(http.MethodPost)
    api.HandleFunc("", put).Methods(http.MethodPut)
    api.HandleFunc("", delete).Methods(http.MethodDelete)

	log.Print("Spinning up the Worst3D Printer Slicing Service...")
    api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)

    

    log.Fatal(http.ListenAndServe(":8080", worstApp.Router))
}
