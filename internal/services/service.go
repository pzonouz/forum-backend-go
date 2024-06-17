package services

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type Service[T any] interface {
	RegisterRoutes()
	GetAll() ([]*T, error)
	GetByID(isTest bool, id int64) (T, error)
	Create(isTest bool, user T) (int64, error)
	EditByID(isTest bool, id int64, user T) error
	DeleteByID(isTest bool, id int64) error
}

func Create[T any](isTest bool, tableName string, instance T, db *sql.DB) (int64, error) {
	var excludedFieldsOfModel []string
	excludedFieldsOfModel = append(excludedFieldsOfModel, "CreatedAt", "ID")

	t := reflect.TypeOf(instance)
	v := reflect.ValueOf(instance)
	query := `INSERT INTO `
	query += `"`

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" `
	query += `(`

	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)

		var j int

		// skip Emty Fields Type add to query
		if v.Field(i).Interface() == "" {
			goto down1
		}
		// skip Ecluded Fields Type add to query
		for j = range excludedFieldsOfModel {
			if t.Field(i).Name == excludedFieldsOfModel[j] {
				goto down1
			}
		}

		query += `"`
		query += sql[0][1]
		query += `"`

		query += `,`
	down1:
	}

	query = query[:len(query)-1]
	query += `)`
	query += ` VALUES `
	query += `(`

	externalI := 0

	for i := 0; i < t.NumField(); i++ {
		var j int
		for j = range excludedFieldsOfModel {
			if t.Field(i).Name == excludedFieldsOfModel[j] {
				goto down
			}
		}

		// query += strconv.Itoa(i)

		//skip Emty Filds to increate query paramrs number
		if v.Field(i).Interface() == "" {
			goto down
		}

		query += `$`
		externalI++
		query += strconv.Itoa(externalI)
		query += `,`
	down:
	}

	query = query[:len(query)-1]
	query += `) RETURNING "id";`
	stmt, err := db.Prepare(query)

	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	return QueryRowWithStruct(stmt, excludedFieldsOfModel, instance, "", true)
}

func Edit[T any](isTest bool, tableName string, db *sql.DB, searchField string, searchFieldValue string, instance T) error {
	var excludedFieldsOfModel []string
	excludedFieldsOfModel = append(excludedFieldsOfModel, "Email", "Password", "CreatedAt", "ID")

	t := reflect.TypeOf(instance)
	v := reflect.ValueOf(instance)
	query := `UPDATE `
	query += `"`

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" `
	query += `SET `
	externalI := 0

	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)

		var j int
		for j = range excludedFieldsOfModel {
			if t.Field(i).Name == excludedFieldsOfModel[j] {
				goto down1
			}
		}

		if v.Field(i).Interface() == "" {
			goto down1
		}

		externalI++
		query += sql[0][1]
		query += `=`
		query += `$`
		query += strconv.Itoa(externalI)
		query += `,`
	down1:
	}

	query = query[:len(query)-1]
	query += ` WHERE `
	query += searchField
	query += `= $`
	query += strconv.Itoa(externalI + 1)

	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = QueryRowWithStruct(stmt, excludedFieldsOfModel, instance, searchFieldValue, false)

	return err
}

func Get[T any](isTest bool, tableName string, db *sql.DB, searchField string, searchFieldValue string, excludedFieldsOfModel []string) (*T, error) {
	query := `SELECT * `
	query += `FROM `
	query += `"`

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" WHERE `
	query += searchField
	query += `= $1`

	stmt, err := db.Prepare(query)

	if err != nil {
		return new(T), err
	}

	defer stmt.Close()

	rows, err := QueryRowsToStruct[T](stmt, excludedFieldsOfModel, searchFieldValue)
	if len(rows) == 0 {
		return new(T), errors.New("not found")
	}

	return rows[0], err
}

func GetAll[T any](isTest bool, tableName string, db *sql.DB, limit string, searchField string, searchFieldValue string, excludedFieldsOfModel []string) ([]*T, error) {
	var objects []*T

	query := `SELECT * `
	query += `FROM `

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	if len(searchField) == 0 {
		query += ` LIMIT `
		query += ` $1`
	} else {
		query += ` LIMIT `
		query += ` $1`
		query += ` WHERE `
		query += searchField
		query += ` $2`
	}

	stmt, err := db.Prepare(query)

	if err != nil {
		return objects, err
	}

	defer stmt.Close()

	var rows []*T

	if len(searchField) != 0 {
		rows, err = QueryRowsToStruct[T](stmt, excludedFieldsOfModel, searchFieldValue, limit)
	} else {
		rows, err = QueryRowsToStruct[T](stmt, excludedFieldsOfModel, limit)
	}

	if len(rows) == 0 {
		return objects, errors.New("not found")
	}

	return rows, err
}

func QueryRowsToStruct[T any](stmt *sql.Stmt, excludedFieldsOfModel []string, args ...any) ([]*T, error) {

	object := new(T)

	var objects []*T

	var rows *sql.Rows

	var err error

	rows, err = stmt.Query(args...)
	if rows.Err() != nil {
		return objects, rows.Err()
	}

	if err != nil {
		return objects, err
	}

	defer rows.Close()

	var params []any

	for rows.Next() {
		params = []any{}
		object = new(T)
		t := reflect.TypeOf(*object)
		v := reflect.ValueOf(object)

		for i := 0; i < t.NumField(); i++ {
			params = append(params, t.Field(i).Name)
			params[i] = &params[i]
		}

		err := rows.Scan(params...)

		if err != nil {
			return objects, err
		}

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			for j := range excludedFieldsOfModel {
				if excludedFieldsOfModel[j] == t.Field(i).Name {
					i++
				}
			}

			paramValue := reflect.ValueOf(params[i])
			if paramValue.Kind() == reflect.Ptr {
				intr := paramValue.Elem().Interface()
				switch intr.(type) {
				case int64:
					v.Field(i).SetInt(intr.(int64))
				case string:
					v.Field(i).SetString(intr.(string))
				}
			} else if paramValue.Kind() == reflect.String {
				v.Field(i).SetString(paramValue.String())
			} else if paramValue.Kind() == reflect.Struct {
				value, ok := paramValue.Interface().(time.Time)
				if !ok {
					return objects, fmt.Errorf("cannot covert to %v", paramValue.Type())
				}

				v.Field(i).SetString(value.String())
			}
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func QueryRowWithStruct[T any](stmt *sql.Stmt, excludedFieldsOfModel []string, instance T, whereClause string, returnResult bool) (int64, error) {
	var id int64

	object := new(T)

	var params []any

	typeOfModel := reflect.TypeOf(*object)
	valueOfModel := reflect.ValueOf(instance)

	var j int

	for i := 0; i < typeOfModel.NumField(); i++ {
		for j = range excludedFieldsOfModel {
			if typeOfModel.Field(i).Name == excludedFieldsOfModel[j] {
				goto down
			}
		}

		if valueOfModel.Field(i).Interface() == "" {
			goto down
		}

		params = append(params, valueOfModel.Field(i).Interface())
	down:
	}

	if len(whereClause) != 0 {
		params = append(params, whereClause)
	}

	var err error
	if returnResult {
		err = stmt.QueryRow(params...).Scan(&id)
	} else {
		_, err = stmt.Exec(params...)
	}

	if err != nil {
		return -1, err
	}

	return id, nil
}
