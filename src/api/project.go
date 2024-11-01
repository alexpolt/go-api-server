package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	// . "go-api-server/common"
	"go-api-server/db"
)

type Project_Id struct {
	User_Id    int64
	Project_Id int64
}

type Project struct {
	User_Id     int64
	Project_Id  int64
	Name        string
	Description string
}

func Load_Projects(ctx context.Context, user_id *int64) (*[]Json_Data, error) {
	query := `SELECT content FROM "main"."projects" WHERE user_id = $1`
	rows, err := db.Pool.Query(ctx, query, *user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects = make([]Json_Data, 0)
	for rows.Next() {
		var project_json Json_Data
		err = rows.Scan(&project_json.data)
		projects = append(projects, project_json)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &projects, err
}

func Get_Project(ctx context.Context, project_id *int64) (*[]byte, error) {
	var project_json []byte
	query := `SELECT content FROM "main"."projects" WHERE project_id = $2`
	row := db.Pool.QueryRow(ctx, query, *project_id)
	err := row.Scan(&project_json)
	return &project_json, err
}

func next_project_id(ctx context.Context) (int64, error) {
	var next_id int64
	query_id := `SELECT nextval('seq.project_ids');`
	row := db.Pool.QueryRow(ctx, query_id)
	err := row.Scan(&next_id)
	return next_id, err
}

func check_project_arguments(project *Project) error {
	if project.Name == "" {
		return fmt.Errorf("name field is empty")
	}
	return nil
}

func Create_Project(ctx context.Context, project *Project) (*[]byte, error) {
	if err := check_project_arguments(project); err != nil {
		return nil, err
	}

	next_id, err := next_project_id(ctx)
	if err != nil {
		return nil, err
	}
	project.Project_Id = next_id
	project_json, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	args := pgx.NamedArgs{
		"user_id":    project.User_Id,
		"project_id": project.Project_Id,
		"content":    project_json,
	}
	query := `INSERT INTO "main"."projects" ( user_id, project_id, content ) ` +
		`VALUES( @user_id, @project_id, @content )`
	_, err = db.Pool.Exec(ctx, query, &args)
	return &project_json, err
}
