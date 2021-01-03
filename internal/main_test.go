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
psql -h localhost -U postgres postgres

Then simply do 'go test -v'

OR- run 'source scripts/run_tests.sh'
*/

package main

import (
	"os"
	"testing"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	"fmt"
	"time"
)

var worst3dps Worst3DPSApp

func TestMain(m *testing.M) {
	user := os.Getenv("WORST_DB_USERNAME")
	if (user == "") {
		user = "postgres"
	}
	worst3dps.Initialize(
		os.Getenv("WORST_DB_USERNAME"),
		os.Getenv("WORST_DB_PASSWORD"),
		os.Getenv("WORST_DB_NAME") )

	ensureProjTableExists()
	code := m.Run()
	clearProjectsTable()
	os.Exit(code)
}

func ensureProjTableExists() {
	if _, err := worst3dps.DB.Exec(projTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// TODO this makes me a bit nervous
func clearProjectsTable() {
	worst3dps.DB.Exec("DELETE FROM projects")
	worst3dps.DB.Exec("ALTER SEQUENCE projects_id_seq RESTART WITH 1")
}


const projTableCreationQuery = `CREATE TABLE IF NOT EXISTS projects
(
    id SERIAL,
    project_name TEXT NOT NULL
);`


// TESTS

func TestHealthCheck(t *testing.T) {
    req, _ := http.NewRequest("GET", "/healthcheck", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code, response)
}

func TestEmptyTable(t *testing.T) {
	clearProjectsTable()

    req, _ := http.NewRequest("GET", "/projects", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code, response)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestGetNonExistentProject(t *testing.T) {
	clearProjectsTable()

    req, _ := http.NewRequest("GET", "/project/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code, response)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Project not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Project not found'. Got '%s'", m["error"])
    }
}

func TestGetProject(t *testing.T) {
	clearProjectsTable()
	addProjects(1)

	req, _ := http.NewRequest("GET", "/project/10000", nil)
	response := executeRequest(req)

	//fmt.Println(response.Body)

	checkResponseCode(t, http.StatusOK, response.Code, response)
}

func TestUpdateProject(t *testing.T) {
	clearProjectsTable()
	addProjects(1)

	req, _ := http.NewRequest("GET", "/project/10000", nil)
	response := executeRequest(req)
	var originalProject map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProject)
	fmt.Println("- GET the existing entry")
	fmt.Println(response.Body)

	var jsonStr = []byte(`{"id":10000,"project_name":"test project with spaces"}`)
	req, _ = http.NewRequest("PUT", "/project/10000", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code, response)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println("- result from PUT")
	fmt.Println(response.Body)

	if m["id"] != originalProject["id"] {
		t.Errorf("Expected the id to not change (%v). Got %v", originalProject["id"], m["id"])
	}

	if m["project_name"] == originalProject["project_name"] {
		t.Errorf("Expected the name to change from %v. Got %v", originalProject["project_name"], m["project_name"])
	}
}

func TestDeleteProject(t *testing.T) {
    clearProjectsTable()
    addProjects(1)  // creates project 10000

	fmt.Println("- check it is there with GET")
    req, _ := http.NewRequest("GET", "/project/10000", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code, response)

	fmt.Println("- do the DELETE")
    req, _ = http.NewRequest("DELETE", "/project/10000", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code, response)

	fmt.Println("- check it is gone with GET")
    req, _ = http.NewRequest("GET", "/project/10000", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code, response)
}

func TestAlwaysWorks(t *testing.T) {
	// do nothing, always happy
}


// generate some projects for testing
func addProjects(count int) {
	for i := 0; i < count; i++ {
		worst3dps.DB.Exec("INSERT INTO projects(id, project_name) VALUES ($1, $2)", i + 10000, "Test_Project_"+strconv.Itoa(i))
	}

	//fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
    time.Sleep(1 * time.Second)
    //fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
}


// lifted this directly from https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    worst3dps.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int, resp *httptest.ResponseRecorder) {
    if expected != actual {
		fmt.Println(resp.Body)
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}
