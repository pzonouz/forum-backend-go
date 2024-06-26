package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"forum-backend-go/internal/middlewares"
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

// DeleteHandler implements Service.
func (u *UserService) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = u.DeleteByID(false, int64(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// EditByID implements Service.
func (u *UserService) EditByID(isTest bool, id int64, user models.User) error {
	return Edit(isTest, "users", u.db, "id", strconv.Itoa(int(id)), user)
}

func (u *UserService) DeleteByID(isTest bool, id int64) error {
	return Delete[models.User](isTest, "users", u.db, "id", strconv.Itoa(int(id)))
}

// GetByID implements Service.
func (u *UserService) GetByID(isTest bool, id int64) (models.User, error) {
	var excludedFields []string
	excludedFields = append(excludedFields, "Password")
	user, err := Get[models.User](isTest, "users", u.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *user, err
	}

	return *user, nil
}

// GetHandler implements Service.
func (u *UserService) GetHandler(w http.ResponseWriter, r *http.Request) {
	access, _ := r.Cookie("access")

	token, _ := jwt.ParseWithClaims(access.Value, &utils.MyClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	claims := token.Claims.(*utils.MyClaims)
	user, err := u.GetByID(false, claims.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	utils.WriteJSON(w, user)
}

// GetHandlerForPlural implements Service.
func (u *UserService) GetHandlerForPlural(w http.ResponseWriter, r *http.Request) {
	var excludedFields []string

	requestQuery := r.URL.Query()
	sortBy := requestQuery.Get("sort_by")
	sortDirection := requestQuery.Get("sort_direction")
	searchField := requestQuery.Get("search_field")
	searchFieldValue := requestQuery.Get("search_field_value")
	operator := requestQuery.Get("operator")
	limit := requestQuery.Get("limit")

	if _, err := strconv.Atoi(limit); err != nil && len(limit) != 0 {
		http.Error(w, "Limit is not number", http.StatusBadRequest)

		return
	}

	if len(sortDirection) != 0 && strings.Compare(sortDirection, "ASC") != 0 && strings.Compare(sortDirection, "DESC") != 0 {
		http.Error(w, "Sort direction in not ASC or DESC", http.StatusBadRequest)

		return
	}

	if len(sortDirection) == 0 {
		sortDirection = "ASC"
	}

	excludedFields = append(excludedFields, "Password", "Role")
	users, err := GetMany[models.User](false, "users", u.db, limit, sortBy, sortDirection, searchField, searchFieldValue, operator, excludedFields)

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.MyClaims{ID: user.ID, Expired: expired.Unix()})
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
		// SameSite: http.SameSiteLaxMode,
		// Domain:   "",
	}
	http.SetCookie(w, &cookie)
	_, _ = w.Write([]byte(signedToken))
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/get_all/", u.GetHandlerForPlural).Methods("GET")
	UsersRouter.HandleFunc("/", middlewares.LoginGuard(u.GetHandler)).Methods("GET")
	UsersRouter.HandleFunc("/register", u.RegisterHandler).Methods("POST")
	UsersRouter.HandleFunc("/login", u.LoginHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", u.PatchHandler).Methods("PATCH")
	UsersRouter.HandleFunc("/{id}", u.DeleteHandler).Methods("DELETE")
}
