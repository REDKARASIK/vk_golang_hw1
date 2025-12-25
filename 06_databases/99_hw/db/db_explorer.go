package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Column struct {
	Name     string
	Type     string
	Nullable bool
	IsPK     bool
}

type Table struct {
	Name    string
	Columns []Column
	PK      string
}

type DBRepository struct {
	db     *sql.DB
	tables map[string]Table
}

func NewDBRepository(db *sql.DB) (*DBRepository, error) {
	repo := &DBRepository{
		db:     db,
		tables: make(map[string]Table),
	}

	err := repo.loadSchema()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

//
// =============== LOAD SCHEMA ===============
//

func (repo *DBRepository) loadSchema() error {
	rows, err := repo.db.Query(`
        SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
        ORDER BY TABLE_NAME, ORDINAL_POSITION
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	current := ""
	var tbl Table

	for rows.Next() {
		var tName, cName, cType, nullable, key string

		if err := rows.Scan(&tName, &cName, &cType, &nullable, &key); err != nil {
			return err
		}

		if tName != current {
			if current != "" {
				repo.tables[current] = tbl
			}
			current = tName
			tbl = Table{Name: tName}
		}

		col := Column{
			Name:     cName,
			Type:     cType,
			Nullable: nullable == "YES",
			IsPK:     key == "PRI",
		}
		if col.IsPK {
			tbl.PK = col.Name
		}
		tbl.Columns = append(tbl.Columns, col)
	}

	if current != "" {
		repo.tables[current] = tbl
	}

	return nil
}

//
// =============== GET / ===============
//

func (repo *DBRepository) handleListTables(w http.ResponseWriter, r *http.Request) {
	var list []string
	for t := range repo.tables {
		list = append(list, t)
	}

	writeResponse(w, map[string]any{"tables": list})
}

//
// =============== GET /table ===============
//

func (repo *DBRepository) handleSelectList(w http.ResponseWriter, r *http.Request, table string) {
	tbl, ok := repo.tables[table]
	if !ok {
		writeError(w, "unknown table", http.StatusNotFound)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf("SELECT * FROM `%s` LIMIT ? OFFSET ?", tbl.Name)
	rows, err := repo.db.Query(query, limit, offset)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	result := []map[string]any{}

	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range ptrs {
			ptrs[i] = &vals[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			writeError(w, err.Error(), 500)
			return
		}

		row := make(map[string]any)
		for i, c := range cols {
			row[c] = convertSQLValue(vals[i])
		}
		result = append(result, row)
	}

	writeResponse(w, map[string]any{"records": result})
}

//
// =============== GET /table/id ===============
//

func (repo *DBRepository) handleSelectOne(w http.ResponseWriter, r *http.Request, table, id string) {
	tbl, ok := repo.tables[table]
	if !ok {
		writeError(w, "unknown table", 404)
		return
	}

	colNames := make([]string, len(tbl.Columns))
	for i, c := range tbl.Columns {
		colNames[i] = c.Name
	}

	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s` = ?", tbl.Name, tbl.PK)
	row := repo.db.QueryRow(query, id)

	ptrs := make([]any, len(colNames))
	vals := make([]any, len(colNames))
	for i := range ptrs {
		ptrs[i] = &vals[i]
	}

	err := row.Scan(ptrs...)
	if err == sql.ErrNoRows {
		writeError(w, "record not found", 404)
		return
	}
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}

	out := map[string]any{}
	for i, c := range colNames {
		out[c] = convertSQLValue(vals[i])
	}

	writeResponse(w, map[string]any{"record": out})
}

//
// =============== PUT /table ===============
//

func (repo *DBRepository) handleInsert(w http.ResponseWriter, r *http.Request, table string) {
	tbl := repo.tables[table]

	var data map[string]any
	json.NewDecoder(r.Body).Decode(&data)

	fields := []string{}
	values := []any{}
	place := []string{}

	for _, col := range tbl.Columns {

		if col.IsPK {
			// PK игнорируется при вставке
			continue
		}

		val, ok := data[col.Name]

		if !ok {
			// Нет значения в body
			if col.Nullable {
				val = nil
			} else {
				// NOT NULL поле → использовать пустую строку
				val = ""
			}
		}

		// проверяем тип
		if val != nil && !validateType(col, val) {
			writeError(w, "field "+col.Name+" have invalid type", 400)
			return
		}

		fields = append(fields, "`"+col.Name+"`")
		values = append(values, val)
		place = append(place, "?")
	}

	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tbl.Name, strings.Join(fields, ", "), strings.Join(place, ", "))
	res, err := repo.db.Exec(query, values...)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	writeResponse(w, map[string]any{tbl.PK: id})
}

//
// =============== POST /table/id ===============
//

func (repo *DBRepository) handleUpdate(w http.ResponseWriter, r *http.Request, table, id string) {
	tbl := repo.tables[table]

	var data map[string]any
	json.NewDecoder(r.Body).Decode(&data)

	set := []string{}
	args := []any{}

	for _, col := range tbl.Columns {
		val, ok := data[col.Name]
		if !ok {
			continue
		}

		if col.IsPK {
			writeError(w, "field "+col.Name+" have invalid type", 400)
			return
		}

		if !validateType(col, val) {
			writeError(w, "field "+col.Name+" have invalid type", 400)
			return
		}

		set = append(set, "`"+col.Name+"` = ?")
		args = append(args, val)
	}

	if len(set) == 0 {
		writeResponse(w, map[string]any{"updated": 0})
		return
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE `%s` = ?", tbl.Name, strings.Join(set, ", "), tbl.PK)
	res, err := repo.db.Exec(query, args...)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}

	n, _ := res.RowsAffected()
	writeResponse(w, map[string]any{"updated": n})
}

//
// =============== DELETE /table/id ===============
//

func (repo *DBRepository) handleDelete(w http.ResponseWriter, r *http.Request, table, id string) {
	tbl := repo.tables[table]

	query := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = ?", tbl.Name, tbl.PK)
	res, err := repo.db.Exec(query, id)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}

	n, _ := res.RowsAffected()
	writeResponse(w, map[string]any{"deleted": n})
}

//
// =============== ROUTER ===============
//

func NewDBExplorer(db *sql.DB) (http.Handler, error) {
	repo, err := NewDBRepository(db)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		if path == "" {
			if r.Method != http.MethodGet {
				writeError(w, "bad method", 400)
				return
			}
			repo.handleListTables(w, r)
			return
		}

		parts := strings.Split(path, "/")

		if len(parts) == 1 {
			table := parts[0]
			switch r.Method {
			case http.MethodGet:
				repo.handleSelectList(w, r, table)
			case http.MethodPut:
				repo.handleInsert(w, r, table)
			default:
				writeError(w, "bad method", 400)
			}
			return
		}

		if len(parts) == 2 {
			table, id := parts[0], parts[1]

			switch r.Method {
			case http.MethodGet:
				repo.handleSelectOne(w, r, table, id)
			case http.MethodPost:
				repo.handleUpdate(w, r, table, id)
			case http.MethodDelete:
				repo.handleDelete(w, r, table, id)
			default:
				writeError(w, "bad method", 400)
			}
			return
		}

		writeError(w, "not found", 404)
	}), nil
}

//
// =============== HELPERS ===============
//

func writeError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{"error": msg})
}

func writeResponse(w http.ResponseWriter, obj any) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"response": obj})
}

func convertSQLValue(v any) any {
	switch val := v.(type) {
	case nil:
		return nil
	case []byte:
		return string(val)
	default:
		return val
	}
}

func validateType(col Column, val any) bool {
	if val == nil {
		return col.Nullable
	}

	switch {
	case strings.HasPrefix(col.Type, "int"):
		_, ok := val.(float64)
		return ok
	case strings.HasPrefix(col.Type, "varchar"), strings.HasPrefix(col.Type, "text"):
		_, ok := val.(string)
		return ok
	}

	return false
}
