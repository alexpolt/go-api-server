package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	// . "go-api-server/common"
	"go-api-server/db"
)

const Namespace_Status_Deleted = "deleted"

type Namespace struct {
	Project_Id   int64
	Namespace_Id int64
	Parent_Id    int64
	Name         string
	Target       string
	Status       string
	Node_Id      int64
	Version      int64
}

type NS_Version struct {
	Namespace_Id int64
	Version      int64
}

func Load_Namespaces(ctx context.Context, project_id *int64) (*[]Json_Data, error) {
	query := `SELECT content || jsonb_build_object('Node_Id', node_id, 'Version', version) FROM "main"."namespaces" WHERE project_id = $1 AND content -> 'Status' != '"%s"' ORDER BY ns_id`
	query = fmt.Sprintf(query, Namespace_Status_Deleted)

	rows, err := db.Pool.Query(ctx, query, *project_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	namespaces := make([]Json_Data, 0)
	for rows.Next() {
		var namespace_json Json_Data
		err = rows.Scan(&namespace_json.data)
		namespaces = append(namespaces, namespace_json)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &namespaces, err
}

func Get_Namespace(ctx context.Context, ns_id *int64) (*[]byte, error) {
	var namespace_json []byte
	query := `SELECT content || jsonb_build_object('Node_Id', node_id, 'Version', version) FROM "main"."namespaces" WHERE ns_id = $1`
	row := db.Pool.QueryRow(ctx, query, ns_id)
	err := row.Scan(&namespace_json)
	return &namespace_json, err
}

func next_namespace_id(ctx context.Context) (int64, error) {
	var next_id int64
	query_id := `SELECT nextval('seq.ns_ids');`
	row := db.Pool.QueryRow(ctx, query_id)
	err := row.Scan(&next_id)
	return next_id, err
}

func check_namespace_arguments(ns *Namespace) error {
	if ns.Project_Id == 0 {
		return fmt.Errorf("project_Id field is 0")
	}
	if ns.Name == "" {
		return fmt.Errorf("name field is empty")
	}
	return nil
}

func Create_Namespace(ctx context.Context, ns *Namespace) (*[]byte, error) {
	if err := check_namespace_arguments(ns); err != nil {
		return nil, err
	}
	next_id, err := next_namespace_id(ctx)
	if err != nil {
		return nil, err
	}
	ns.Namespace_Id = next_id
	namespace_json, err := json.Marshal(ns)
	if err != nil {
		return nil, err
	}
	args := pgx.NamedArgs{
		"project_id":   ns.Project_Id,
		"namespace_id": ns.Namespace_Id,
		"content":      namespace_json,
	}
	query := `INSERT INTO "main"."namespaces" ( project_id, ns_id, content ) ` +
		`VALUES( @project_id, @namespace_id, @content )`
	_, err = db.Pool.Exec(ctx, query, &args)
	return &namespace_json, err
}

func Delete_Namespace(ctx context.Context, ns_id *int64) (*Call_Success, error) {
	query := `UPDATE "main"."namespaces" SET Content['Status'] = '"%s"' WHERE ns_id = $1`
	query = fmt.Sprintf(query, Namespace_Status_Deleted)
	_, err := db.Pool.Exec(ctx, query, ns_id)
	if err != nil {
		return nil, err
	}
	call_success := &Call_Success{true}
	return call_success, nil
}
