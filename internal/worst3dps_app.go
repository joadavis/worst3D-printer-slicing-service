// worst3dps_app.go
// joadavis Nov 2020

package main

import (
	"database/sql"
    "fmt"
    "log"

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
}

func (app *Worst3DPSApp) Run(addr string) { }