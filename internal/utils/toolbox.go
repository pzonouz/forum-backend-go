package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ReadJSON[T any](w http.ResponseWriter, r *http.Request) T {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}

	decoder := json.NewDecoder(r.Body)
	data := new(T)
	err := decoder.Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return *data
}

func WriteJSON(w http.ResponseWriter, data any) {
	JSON := json.NewEncoder(w)
	err := JSON.Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

type MyClaims struct {
	ID          int64
	Email       string
	Address     string
	NickName    string
	PhoneNumber string
	Role        string
	Expired     int64
	jwt.RegisteredClaims
}

func GetManyQueryCreator(isTest bool, tableName string, sortBy string, operator string, searchField string, sortDirection string) string {
	query := `SELECT * `
	query += `FROM `

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	if sortBy == "" {
		sortBy = "created_at"
	}

	if operator == "" {
		operator = "="
	}

	if len(searchField) == 0 {
		query += ` ORDER BY `
		query += sortBy
		query += ` `
		query += sortDirection
		query += ` LIMIT`
		query += ` $1`
	} else {
		query += ` WHERE `
		query += searchField
		query += ` `
		query += operator
		query += ` $1`
		query += ` ORDER BY `
		query += sortBy
		query += ` `
		query += sortDirection

		query += ` LIMIT `
		query += ` $2`
	}

	return query
}

func GetQueryCreator(isTest bool, tableName string, searchField string) string {
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

	return query
}

func DeleteQueryCreator(isTest bool, tableName string, searchField string) string {
	query := `DELETE FROM `
	query += `"`

	if isTest {
		query += tableName + `_test`
	} else {
		query += tableName
	}

	query += `" `
	query += ` WHERE `
	query += searchField
	query += `= $1`

	return query
}

func GetUserFromRequest(r *http.Request, w http.ResponseWriter) (*MyClaims, error) {
	access, err := r.Cookie("access")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		access.Value,
		&MyClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(*MyClaims)

	if claims.Expired < time.Now().Unix() {
		http.Error(w, "expired", http.StatusUnauthorized)
	}

	return claims, nil
}

func GetUserRoleFromRequest(r *http.Request, w http.ResponseWriter) string {
	access, _ := r.Cookie("access")

	token, err := jwt.ParseWithClaims(
		access.Value,
		&MyClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)
	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(*MyClaims)
	if claims.Expired < time.Now().Unix() {
		http.Error(w, "expired", http.StatusUnauthorized)
	}

	return claims.Role
}

func GetScoreOfUserToQuestion(db *sql.DB, user_id int64, question_id int64) (int64, error) {
	query := `SELECT COALESCE(SUM(CASE
             WHEN operator = 'plus' THEN 1
              ELSE -1
             END),0) AS total
            FROM scores WHERE user_id=$1 AND question_id=$2;`
	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var result int64
	err = stmt.QueryRow(user_id, question_id).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetScoreOfUserToAnswer(db *sql.DB, user_id int64, answer_id int64) (int64, error) {
	query := `SELECT COALESCE(SUM(CASE
             WHEN operator = 'plus' THEN 1
              ELSE -1
             END),0) AS total
            FROM scores WHERE user_id=$1 AND answer_id=$2;`
	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var result int64
	err = stmt.QueryRow(user_id, answer_id).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetScoreOfQuestion(db *sql.DB, question_id int64) (int64, error) {
	query := `SELECT COALESCE(SUM(CASE
             WHEN operator = 'plus' THEN 1
              ELSE -1
             END),0) AS total
            FROM scores WHERE question_id=$1;`
	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var result int64
	err = stmt.QueryRow(question_id).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetScoreOfAnswer(db *sql.DB, answer_id int64) (int64, error) {
	query := `SELECT COALESCE(SUM(CASE
             WHEN operator = 'plus' THEN 1
              ELSE -1
             END),0) AS total
            FROM scores WHERE answer_id=$1;`
	stmt, err := db.Prepare(query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var result int64
	err = stmt.QueryRow(answer_id).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func ResetScoreOfUserToQustion(db *sql.DB, user_id int64, question_id int64) error {
	query := `DELETE FROM scores WHERE user_id=$1 AND question_id=$2`
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user_id, question_id)
	if err != nil {
		return err
	}

	return nil
}

func ResetScoreOfUserToAnswer(db *sql.DB, user_id int64, answer_id int64) error {
	query := `DELETE FROM scores WHERE user_id=$1 AND answer_id=$2`
	stmt, err := db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user_id, answer_id)
	if err != nil {
		return err
	}

	return nil
}
