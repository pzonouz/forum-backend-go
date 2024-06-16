package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

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
	return EditByID(isTest, "users", u.db, id, user)
}

// GetAll implements Service.
func (u *UserService) GetAll() ([]*models.User, error) {
	var excludedFields []string
	excludedFields = append(excludedFields, "id")
	users, err := GetAll[models.User](false, "users", u.db, "20", "", "", excludedFields)

	if err != nil {
		return users, err
	}

	return users, nil
}

// GetByID implements Service.
func (u *UserService) GetByID(isTest bool, id int64) (models.User, error) {
	var excludedFields []string
	excludedFields = append(excludedFields, "id")
	user, err := Get[models.User](isTest, "users", u.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *user, err
	}

	return *user, nil
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
func (u *UserService) GetHandlerForPlural(w http.ResponseWriter, _ *http.Request) {
	users, err := u.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, users)
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
func (u *UserService) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ReadJSON[models.User](w, r)
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 3)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user.Password = string(encryptedPassword)
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

func (u *UserService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	userJSON := utils.ReadJSON[models.User](w, r)

	var excludedFields []string
	user, err := Get[models.User](false, "users", u.db, "email", userJSON.Email, excludedFields)

	if err != nil {
		http.Error(w, "Login and Password does not match", http.StatusUnauthorized)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userJSON.Password))
	if err != nil {
		http.Error(w, "Login and Password does not match", http.StatusUnauthorized)

		return
	}
	expired := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": user.ID, "expired": expired.Unix()})
	signedToken, err := token.SignedString([]byte("secret"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	cookie := http.Cookie{
		Path:     "/",
		Name:     "access",
		Value:    signedToken,
		Expires:  expired,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	_, _ = w.Write([]byte(signedToken))
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	APIV1Router.HandleFunc("", u.GetHandlerForPlural)
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/", u.GetHandlerForPlural).Methods("GET")
	UsersRouter.HandleFunc("/{id}", u.GetHandler).Methods("GET")
	UsersRouter.HandleFunc("/register", u.RegisterHandler).Methods("POST")
	UsersRouter.HandleFunc("/login", u.LoginHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", u.PatchHandler).Methods("PATCH")
	UsersRouter.HandleFunc("/{id}", u.DeleteHandler).Methods("DELETE")
}
