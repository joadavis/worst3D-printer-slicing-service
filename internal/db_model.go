// db_model.go
// Database model to represent users, projects, and jobs
// joadavis Nov 2020

package main

import (
	"database/sql"
	"errors"
)

type user struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Project_ID int    `json:"project_id"`
}

type project struct {
	ID   int `json:"id"`
	Name string `json:"project_name"`
}


func (p *project) getProject(db *sql.DB) error {
	return errors.New("Not implemented")
}