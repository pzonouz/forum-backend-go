package services

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	user, err := Get[models.User](
		isTest,
		"users",
		u.db,
		"id",
		strconv.Itoa(int(id)),
		excludedFields,
	)

	if err != nil {
		return *user, err
	}

	return *user, nil
}

// GetHandler implements Service.
func (u *UserService) GetHandler(w http.ResponseWriter, r *http.Request) {
	requestUser, err := utils.GetUserFromRequest(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := u.GetByID(false, requestUser.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	utils.WriteJSON(w, user)
}

func (u *UserService) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := u.GetByID(false, int64(ID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
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

	if len(sortDirection) != 0 && strings.Compare(sortDirection, "ASC") != 0 &&
		strings.Compare(sortDirection, "DESC") != 0 {
		http.Error(w, "Sort direction in not ASC or DESC", http.StatusBadRequest)

		return
	}

	if len(sortDirection) == 0 {
		sortDirection = "ASC"
	}

	excludedFields = append(excludedFields, "Password")
	users, err := GetMany[models.User](
		false,
		"users",
		u.db,
		limit,
		sortBy,
		sortDirection,
		searchField,
		searchFieldValue,
		operator,
		excludedFields,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, users)
}

// PatchHandler implements Service.
func (u *UserService) PatchHandlerAdmin(w http.ResponseWriter, r *http.Request) {
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

func (u *UserService) PatchHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, err := utils.GetUserFromRequest(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return
	}

	user := utils.ReadJSON[models.User](w, r)

	err = u.EditByID(false, currentUser.ID, user)

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

func (u *UserService) IsUniqueEmailHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ReadJSON[models.User](w, r)
	userFromReqest, _ := utils.GetUserFromRequest(r, w)
	currentUser, _ := Get[models.User](false, "users", u.db, "email", user.Email, nil)
	if len(currentUser.Email) > 0 && currentUser.Email != userFromReqest.Email {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func (u *UserService) IsUniqueNickNameHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ReadJSON[models.User](w, r)
	userFromReqest, _ := utils.GetUserFromRequest(r, w)
	currentUser, _ := Get[models.User](false, "users", u.db, "nickName", user.NickName, nil)
	if len(currentUser.NickName) > 0 && currentUser.NickName != userFromReqest.NickName {
		http.Error(w, "", http.StatusBadRequest)
	}

}

func (u *UserService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	userJSON := utils.ReadJSON[models.User](w, r)

	print(userJSON.Password)
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
	expired := time.Now().Add(time.Hour * 24 * 365 * 5)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		utils.MyClaims{ID: user.ID, Expired: expired.Unix(), Role: user.Role, NickName: user.NickName, Email: user.Email},
	)
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

	type data struct {
		Access string `json:"access"`
	}

	utils.WriteJSON(w, &data{Access: signedToken})
}

func (u *UserService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Path:     "/",
		Name:     "access",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		// SameSite: http.SameSiteLaxMode,
		// Domain:   "",
	}
	http.SetCookie(w, &cookie)
}

func (u *UserService) GetGoogleOauthLinkHandler(w http.ResponseWriter, r *http.Request) {
	link := `https://accounts.google.com/o/oauth2/v2/auth?scope=https://www.googleapis.com/auth/userinfo.email+https://www.googleapis.com/auth/userinfo.profile+openid&access_type=offline&include_granted_scopes=true&response_type=code&redirect_uri=https%3A//localhost/api/v1/users/google-callback&client_id=540094082819-cvffsbg31rcsva57f0fne7urqt34d6ur.apps.googleusercontent.com`

	type data struct {
		Link string `json:"link"`
	}

	utils.WriteJSON(w, &data{Link: link})
}

func (u *UserService) GooleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	authorizationExchangeUrl := "https://oauth2.googleapis.com/token"
	resp, err := http.PostForm(authorizationExchangeUrl, url.Values{"client_id": {"540094082819-cvffsbg31rcsva57f0fne7urqt34d6ur.apps.googleusercontent.com"}, "client_secret": {"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"}, "code": {code}, "grant_type": {"authorization_code"}, "redirect_uri": {"https%3A//localhost/api/v1/users/google-callback"}})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer resp.Body.Close()

	log.Printf("%v", string(data))

	type output struct {
		Data string `json:data`
	}
	utils.WriteJSON(w, &output{Data: string(data)})
}

func (u *UserService) ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	token := uuid.New()
	tokenString := strings.ReplaceAll(string(token.String()), "-", "")
	query := `UPDATE users SET is_forget_password=true,token='` + tokenString + `' WHERE email=$1`
	result, err := u.db.Exec(query, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No Email", http.StatusBadRequest)
		return
	}
	gmailUsername := utils.GetEnv("GMAIL_USERNAME", "p.zonouz@gmail.com")
	gmailPassword := utils.GetEnv("GMAIL_PASSWORD", "egfu usxu dcmv xblu")
	gmailPort := utils.GetEnv("GMAIL_PORT", "587")
	gmailAddress := utils.GetEnv("GMAIL_SMTP_ADDRESS", "smtp.gmail.com")
	gmailAuth := smtp.PlainAuth("", gmailUsername, gmailPassword, gmailAddress)
	from := utils.GetEnv("GMAIL_FROM", "ICEF")
	to := []string{email}
	t, err := template.ParseFiles("../../internal/services/template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: ایمیل تغییر پسورد \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Address string
	}{
		Address: `https://localhost/users/forget_password_callback/` + tokenString})

	err = smtp.SendMail(gmailAddress+":"+gmailPort, gmailAuth, from, to, body.Bytes())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (u *UserService) ForgetPasswordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	data := utils.ReadJSON[models.User](w, r)
	if len(data.Password) < 6 {
		http.Error(w, "Password length", http.StatusBadRequest)
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 3)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	query := `UPDATE users SET password=$1 WHERE token=$2 AND is_forget_password=true`
	tx, err := u.db.BeginTx(context.Background(), nil)
	defer tx.Rollback()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	result, err := tx.ExecContext(context.Background(), query, data.Password, token)
	if err == nil {
		rowsCount, _ := result.RowsAffected()
		if rowsCount == 0 {
			http.Error(w, "Token is wrong", http.StatusBadRequest)
			return
		}
	}
	query = `UPDATE users SET is_forget_password=false,password=$1 WHERE token=$2`
	result, err = tx.ExecContext(context.Background(), query, encryptedPassword, token)
	if err == nil {
		rowsCount, _ := result.RowsAffected()
		if rowsCount == 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/get_all/", u.GetHandlerForPlural).Methods("GET")
	UsersRouter.HandleFunc("/", middlewares.LoginGuard(u.GetHandler)).Methods("GET")
	UsersRouter.HandleFunc("/get_by_id/{id}", u.GetByIDHandler).Methods("GET")
	UsersRouter.HandleFunc("/google-callback", u.GooleCallbackHandler).Methods("GET")
	UsersRouter.HandleFunc("/register", u.RegisterHandler).Methods("POST")
	UsersRouter.HandleFunc("/is_unique_email", u.IsUniqueEmailHandler).Methods("POST")
	UsersRouter.HandleFunc("/is_unique_nickname", u.IsUniqueNickNameHandler).Methods("POST")
	UsersRouter.HandleFunc("/login", u.LoginHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", middlewares.AdminRoleGuard(u.PatchHandler)).Methods("PATCH")
	UsersRouter.HandleFunc("/", middlewares.LoginGuard(u.PatchHandler)).Methods("PATCH")
	UsersRouter.HandleFunc("/logout", middlewares.LoginGuard(u.LogoutHandler)).Methods("GET")
	UsersRouter.HandleFunc("/get_google_oauth_link", u.GetGoogleOauthLinkHandler).Methods("GET")
	UsersRouter.HandleFunc("/{id}", middlewares.AdminRoleGuard(u.DeleteHandler)).Methods("DELETE")
	UsersRouter.HandleFunc("/forget_password/{email}", u.ForgetPasswordHandler).Methods("GET")
	UsersRouter.HandleFunc("/forget_password_callback/{token}/", u.ForgetPasswordCallbackHandler).Methods("POST")
}
