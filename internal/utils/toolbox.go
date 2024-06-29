package utils

import (
	"encoding/json"
	"log"
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

func GetUserIDFromRequest(r *http.Request, w http.ResponseWriter) int64 {
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
	return claims.ID
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
	log.Print(claims.Role)
	if claims.Expired < time.Now().Unix() {
		http.Error(w, "expired", http.StatusUnauthorized)
	}

	return claims.Role
}
