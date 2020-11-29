// main_test.go
// joadavis Nov 2020

/* 
Please set up environment first
export WORST_DB_USERNAME=postgres
export WORST_DB_PASSWORD=
export WORST_DB_NAME=postgres

Then spin up a postgres instance with
sudo docker run -it --name test-pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

And manually test connection with
psql -h localhost -u postgres postgres
*/

package main

import (
	"os"
	"testing"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

var worst3dps Worst3DPSApp

func TestMain(m *testing.M) {
	worst3dps.Initialize(
		os.Getenv("WORST_DB_USERNAME"),
		os.Getenv("WORST_DB_PASSWORD"),
		os.Getenv("WORST_DB_NAME") )

	ensureProjTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureProjTableExists() {
	if _, err := worst3dps.DB.Exec(projTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// TODO this makes me a bit nervous
func clearTable() {
	worst3dps.DB.Exec("DELETE FROM projects")
	worst3dps.DB.Exec("ALTER SEQUENCE projects_id_seq RESTART WITH 1")
}


const projTableCreationQuery = `CREATE TABLE IF NOT EXISTS projects
(
    id SERIAL,
    project_name TEXT NOT NULL
);`


// TESTS

func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/projects", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestGetNonExistentProject(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/project/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Project not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Project not found'. Got '%s'", m["error"])
    }
}





// lifted this directly from https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    worst3dps.Router.ServeHTTP(rr, req)

    return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}