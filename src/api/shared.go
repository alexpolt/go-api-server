package api

import "github.com/jackc/pgx/v5"

// . "go-api-server/common"

type Call_Success struct {
	Success bool
}

type Null_Argument struct{}

type Json_Data struct {
	data []byte
}

func (doc *Json_Data) MarshalJSON() ([]byte, error) {
	return doc.data, nil
}

func (doc *Json_Data) UnmarshalJSON(data []byte) error {
	doc.data = data
	return nil
}

func Rows_To_Json_Data(rows pgx.Rows) (*[]Json_Data, error) {
	var err error
	nodes := make([]Json_Data, 0)
	for rows.Next() {
		var node_json Json_Data
		err = rows.Scan(&node_json.data)
		nodes = append(nodes, node_json)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &nodes, err
}
