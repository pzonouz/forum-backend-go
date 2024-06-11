package services

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/models"
	"forum-backend-go/internal/utils"
)

func NewUserService(db *sql.DB, router *mux.Router) *UserService {
	return &UserService{
		db: db, router: router,
	}
}

type UserService struct {
	db     *sql.DB
	router *mux.Router
}

// Create implements Service.
func (u *UserService) Create(isTest bool, user models.User) (int64, error) {
	var id int64
	id, err := Create[models.User](isTest, "users", user, u.db)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (u *UserService) DeleteByID(isTest bool, id int64) error {
	var stmt *sql.Stmt

	var err error

	if isTest {
		stmt, err = u.db.Prepare(utils.DeleteUserByIDQueryTest)
		if err != nil {
			return err
		}
	} else {
		stmt, err = u.db.Prepare(utils.DeleteUserByIDQuery)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

// DeleteHandler implements Service.
func (u *UserService) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// EditByID implements Service.
func (u *UserService) EditByID(isTest bool, id int64, user models.User) error {
	var stmt *sql.Stmt

	var err error

	if isTest {
		stmt, err = u.db.Prepare(utils.EditUserQueryTest)
	} else {
		stmt, err = u.db.Prepare(utils.EditUserQuery)
	}

	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Address, user.PhoneNumber, id)
	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

// GetAll implements Service.
func (u *UserService) GetAll() ([]models.User, error) {
	panic("unimplemented")
}

// GetByID implements Service.
func (u *UserService) GetByID(isTest bool, id int64) (models.User, error) {
	var excluded_fields []string
	excluded_fields = append(excluded_fields, "id")
	user, err := Get[models.User](isTest, "users", u.db, "id", strconv.Itoa(int(id)), excluded_fields)
	if err != nil {
		return *user, err
	}
	return *user, nil
	// var user models.User

	// var stmt *sql.Stmt

	// var err error

	// if isTest {
	// 	stmt, err = u.db.Prepare(`SELECT id,email,first_name,last_name,address,phone_number,created_at from users_test where id=$1`)
	// } else {
	// 	stmt, err = u.db.Prepare(`SELECT id,email,first_name,last_name,address,phone_number,created_at from users where id=$1`)
	// }

	// if err != nil {
	// 	return user, err
	// }

	// err = stmt.QueryRow(id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Address, &user.PhoneNumber, &user.CreatedAt)

	// defer stmt.Close()

	// return user, err
}

// GetHandler implements Service.
func (u *UserService) GetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := u.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, user)
}

// GetHandlerForPlural implements Service.
func (u *UserService) GetHandlerForPlural(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Users"))
}

// PatchHandler implements Service.
func (u *UserService) PatchHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := utils.ReadJSON[models.User](w, r)
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = u.EditByID(false, int64(ID), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// PostHandler implements Service.
func (u *UserService) PostHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ReadJSON[models.User](w, r)
	id, err := u.Create(false, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	APIV1Router.HandleFunc("", u.GetHandlerForPlural)
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/", u.GetHandlerForPlural).Methods("GET")
	UsersRouter.HandleFunc("/{id}", u.GetHandler).Methods("GET")
	UsersRouter.HandleFunc("/register", u.PostHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", u.PatchHandler).Methods("PATCH")
	UsersRouter.HandleFunc("/{id}", u.DeleteHandler).Methods("DELETE")
}
