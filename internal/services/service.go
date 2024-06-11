package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type Service[T any] interface {
	RegisterRoutes()
	GetHandler(w http.ResponseWriter, r *http.Request)
	GetHandlerForPlural(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
	PatchHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	GetAll() ([]T, error)
	GetByID(isTest bool, id int64) (T, error)
	Create(isTest bool, user T) (int64, error)
	EditByID(isTest bool, id int64, user T) error
	DeleteByID(isTest bool, id int64) error
}

func Create[T any](isTest bool, tableName string, instance T, db *sql.DB, excluded_fields []string) (int64, error) {
	var id int64

	excluded_fields = append(excluded_fields, "id", "created_at")
	query := `INSERT INTO `
	query += `"`
	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" `
	query += `(`
	t := reflect.TypeOf(&instance)

	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
	} else {
		return -1, errors.New("must be struct")
	}

	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)

		for j := range excluded_fields {
			if sql[0][1] == excluded_fields[j] {
				goto down
			}
		}

		query += `"`
		query += sql[0][1]
		query += `"`
		query += `,`
	down:
	}

	query = query[:len(query)-1]
	query += `)`
	query += ` VALUES `
	query += `(`

	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)

		if (sql[0][1] == "id") || (sql[0][1] == "created_at") {
			continue
		}

		query += `$`
		query += strconv.Itoa(i)
		query += `,`
	}

	query = query[:len(query)-1]
	query += `) RETURNING "id";`
	smtmt, err := db.Prepare(query)

	if err != nil {
		return -1, err
	}

	defer smtmt.Close()

	var params []interface{}

	t2 := reflect.TypeOf(&instance)
	v2 := reflect.ValueOf(&instance)

	if t2.Kind() == reflect.Ptr || t2.Elem().Kind() == reflect.Struct {
		t2 = t2.Elem()
		v2 = v2.Elem()
	} else {
		return -1, errors.New("instance is not struct")
	}

	for i := 0; i < t2.NumField(); i++ {
		for j := range excluded_fields {
			re := regexp.MustCompile(`sql:\"(\w+)\"`)
			sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)

			if sql[0][1] == excluded_fields[j] {
				goto down2
			}
		}

		params = append(params, v2.Field(i).String())
	down2:
	}

	err = smtmt.QueryRow(params...).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func Edit[T any](isTest bool, tableName string, instance T, db *sql.DB, excluded_fields []string) error {

	excluded_fields = append(excluded_fields, "id", "created_at")
	query := `INSERT INTO `
	query += `"`

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" `
	query += `(`
	t := reflect.TypeOf(&instance)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
	} else {
		return errors.New("must be struct")
	}
	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)
		for j := range excluded_fields {
			if sql[0][1] == excluded_fields[j] {
				goto down
			}
		}
		query += `"`
		query += sql[0][1]
		query += `"`
		query += `,`
	down:
	}
	query = query[:len(query)-1]
	query += `)`
	query += ` VALUES `
	query += `(`
	for i := 0; i < t.NumField(); i++ {
		re := regexp.MustCompile(`sql:\"(\w+)\"`)
		sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)
		if (sql[0][1] == "id") || (sql[0][1] == "created_at") {
			continue
		}
		query += `$`
		query += strconv.Itoa(i)
		query += `,`
	}

	query = query[:len(query)-1]
	query += `) RETURNING "id";`
	smtmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	var params []interface{}

	t2 := reflect.TypeOf(&instance)
	v2 := reflect.ValueOf(&instance)

	if t2.Kind() == reflect.Ptr || t2.Elem().Kind() == reflect.Struct {
		t2 = t2.Elem()
		v2 = v2.Elem()
	} else {
		return errors.New("instance is not struct")
	}

	for i := 0; i < t2.NumField(); i++ {
		for j := range excluded_fields {
			re := regexp.MustCompile(`sql:\"(\w+)\"`)
			sql := re.FindAllStringSubmatch(string(t.Field(i).Tag), 1)
			if sql[0][1] == excluded_fields[j] {
				goto down2
			}
		}

		params = append(params, v2.Field(i).String())
	down2:
	}

	row := smtmt.QueryRow(params...)
	err = row.Err()

	if row.Err() != nil {
		return err
	}

	return err
}

func Get[T any](isTest bool, tableName string, db *sql.DB, searchField string, searchFieldValue string, excludedFieldsOfModel []string) (*T, error) {
	excludedFieldsOfModel = append(excludedFieldsOfModel, "Password")
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

	return QueryRowToStruct[T](stmt, excludedFieldsOfModel, searchFieldValue)
}

func QueryRowToStruct[T any](stmt *sql.Stmt, excludedFieldsOfModel []string, args ...any) (*T, error) {
	var params []any

	object := new(T)

	t := reflect.TypeOf(*object)
	v := reflect.ValueOf(object)

	for i := 0; i < t.NumField(); i++ {
		params = append(params, t.Field(i).Name)
		params[i] = &params[i]
	}

	err := stmt.QueryRow(args...).Scan(params...)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if err != nil {
		return object, err
	}

	for i := 0; i < t.NumField(); i++ {
		var j int
		for j, _ = range excludedFieldsOfModel {
			if excludedFieldsOfModel[j] == t.Field(i).Name {
				fmt.Printf("%v-%v", i, t.Field(i).Name)
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
				return object, fmt.Errorf("cannot covert to %v", paramValue.Type())
			}

			v.Field(i).SetString(value.String())
		}
	}

	return object, nil
}
