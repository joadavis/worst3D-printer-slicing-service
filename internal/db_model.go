// db_model.go
// Database model to represent users, projects, and jobs
// and roles to tie them together
// joadavis Nov 2020

package main

// TODO manually remove errors
import (
	"database/sql"
	"errors"
)

type user struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type project struct {
	ID   int `json:"id"`
	Name string `json:"project_name"`
}

type job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}



// Projects
func (p *project) createProject(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO projects(id, project_name) VALUES ($1, $2) RETURNING id",
		p.ID, p.Name).Scan(&p.ID)
	
	if err != nil {
		return err
	}

	return nil
}

func (p *project) getProject(db *sql.DB) error {
	return db.QueryRow("SELECT id, project_name FROM projects where id=$1", p.ID).Scan(&p.ID, &p.Name)
}

func (p *project) updateProject(db *sql.DB) error {
	_, err := db.Exec("UPDATE projects SET project_name=$1 where id=$2", p.Name, p.ID)

	return err
}

// TODO: Consider not allowing this function in 'production' as it could have consequences
func (p *project) deleteProject(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM projects WHERE id=$1", p.ID)

	return err
}

// plural action
func getProjects(db *sql.DB, start, count int) ([]project, error) {
	rows, err := db.Query(
		"SELECT id, project_name FROM projects LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []project{}

	for rows.Next() {
		var p project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}


func (u *user) getUser(db *sql.DB) error {
	return errors.New("Not implemented")
}
