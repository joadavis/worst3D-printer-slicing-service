// worst3dps_app.go
// joadavis Nov 2020

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"strconv"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Worst3DPSApp struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *Worst3DPSApp) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
    app.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

	app.Router = mux.NewRouter()
	
	app.initializeRoutes()
}

func (app *Worst3DPSApp) initializeRoutes() {
	app.Router.HandleFunc("/projects", app.getProjects).Methods("GET")
	//app.Router.HandleFunc("/project", app.createProject).Methods("POST")
	app.Router.HandleFunc("/project/{id:[0-9]*}", app.updateProject).Methods(http.MethodPut)
	app.Router.HandleFunc("/project/{id:[0-9]+}", app.getProject).Methods(http.MethodGet)
	app.Router.HandleFunc("/project/{id:[0-9]*}", app.deleteProject).Methods(http.MethodDelete)

	// users

	// token 

	// jobs
	
	// health check
	app.Router.HandleFunc("/healthcheck", app.checkHealth).Methods(http.MethodGet)
}

func (app *Worst3DPSApp) Run(addr string) { }


func (app *Worst3DPSApp) checkHealth(w http.ResponseWriter, r *http.Request) {
	// TODO what else is worth reporting?

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

// Projects
func (app *Worst3DPSApp) getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}
	// TODO authorize against the user and project
	// TODO look up how to do an http unauthorized

	p:= project{ID: id}
	if err := p.getProject(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Project not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (app *Worst3DPSApp) updateProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid project ID")
        return
    }

    var p project
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
		fmt.Println(err)
		fmt.Println(r.Body)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload: ")
        return
    }
    defer r.Body.Close()
    p.ID = id

    if err := p.updateProject(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (app *Worst3DPSApp) deleteProject(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Project ID")
        return
    }

    p := project{ID: id}
    if err := p.deleteProject(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// multiple projects
func (app *Worst3DPSApp) getProjects(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    products, err := getProjects(app.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, products)
}



// Lifted these generic functions from https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}