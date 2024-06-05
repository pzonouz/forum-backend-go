package services

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"forum-backend-go/ineternal/models"
	"forum-backend-go/ineternal/utils"
)

type Service[T any] interface {
	RegisterRoutes()
	GetHandler(w http.ResponseWriter, r *http.Request)
	GetHandlerForPlurar(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
	PatchHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	GetAll() ([]T, error)
	GetByID(isTest bool, id int64) (T, error)
	Create(isTest bool, user T) (int64, error)
	EditByID(isTest bool, id int64, user T) error
	DeleteByID(isTest bool, id int64) error
}

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

	var smtmt *sql.Stmt

	var err error

	if isTest {
		smtmt, err = u.db.Prepare(utils.CreateUserQueryTest)
	} else {
		smtmt, err = u.db.Prepare(utils.CreateUserQuery)
	}

	if err != nil {
		return 0, err
	}

	err = smtmt.QueryRow(user.Email, user.Password, user.FirstName, user.LastName, user.Address, user.PhoneNumber).Scan(&id)

	defer smtmt.Close()

	if err != nil {
		return 0, err
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
	stmt, err = u.db.Prepare(utils.EditUserQueryTest)

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
	var user models.User

	var stmt *sql.Stmt

	var err error

	if isTest {
		stmt, err = u.db.Prepare(`SELECT * from users_test where id=$1`)
	} else {
		stmt, err = u.db.Prepare(`SELECT * from users where id=$1`)
	}

	if err != nil {
		return user, err
	}
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Address, &user.PhoneNumber, &user.CreatedAt)

	defer stmt.Close()

	return user, err
}

// GetHandler implements Service.
func (u *UserService) GetHandler(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
}

// GetHandlerForPlurar implements Service.
func (u *UserService) GetHandlerForPlurar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Users"))
}

// PatchHandler implements Service.
func (u *UserService) PatchHandler(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
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
	APIV1Router.HandleFunc("", u.GetHandlerForPlurar)
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/", u.GetHandlerForPlurar).Methods("GET")
	UsersRouter.HandleFunc("/{id}", u.GetHandler).Methods("GET")
	UsersRouter.HandleFunc("/register", u.PostHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", u.PatchHandler).Methods("PATCH")
	UsersRouter.HandleFunc("/{id}", u.DeleteHandler).Methods("DELETE")
}
