package services

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/middlewares"
	"forum-backend-go/internal/models"
	"forum-backend-go/internal/utils"
)

func NewFileRequestService(db *sql.DB, router *mux.Router) *FileRequest {
	return &FileRequest{
		db:     db,
		router: router,
	}
}

type FileRequest struct {
	db     *sql.DB
	router *mux.Router
}

func (r *FileRequest) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
	searchField := req.URL.Query().Get("search_field")
	searchFieldValue := req.URL.Query().Get("search_field_value")
	operator := req.URL.Query().Get("operator")
	files, err := GetMany[models.File](false, "files", r.db, "", "", "", searchField, searchFieldValue, operator, []string{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, files)
}

func (r *FileRequest) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	file, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, file)
}

func (r *FileRequest) PostHandler(w http.ResponseWriter, req *http.Request) {
	file := utils.ReadJSON[models.File](w, req)
	id, err := r.Create(false, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (r *FileRequest) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}

	filePartial := utils.ReadJSON[models.File](w, req)
	_, err = r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, "not Found", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)

	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if strings.Compare(user.Role, "admin") != 0 {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if strings.Compare(filePartial.Title, "") != 0 && len(filePartial.Title) < 5 {
		http.Error(w, "At least 5 character for name", http.StatusBadRequest)

		return
	}

	err = r.EditByID(false, int64(id), filePartial)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *FileRequest) DeleteAdminHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}

	file, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, "not Found", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)

	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if strings.Compare(user.Role, "admin") != 0 {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}
	err = r.DeleteByID(false, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = os.Remove("../../uploads/" + file.FileName)
}

func (r *FileRequest) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}

	file, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, "not Found", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)

	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if strings.Compare(user.Role, "admin") != 0 && file.UserID != user.ID {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}
	if strings.Compare(file.Title, "") != 0 {
		query := `UPDATE files SET question_id=NULL,answer_id=NULL WHERE id=$1`
		result, err := r.db.Exec(query, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if rows == 0 {
			http.Error(w, "0 rows affected", http.StatusBadRequest)
			return
		}
		return
	}
	_ = os.Remove("../../uploads/" + file.FileName)
	err = r.DeleteByID(false, int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// Create implements Service.
func (r *FileRequest) Create(isTest bool, file models.File) (int64, error) {
	var id int64
	id, err := Create[models.File](isTest, "files", file, r.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (r *FileRequest) DeleteByID(isTest bool, id int64) error {
	return Delete[models.File](isTest, "files", r.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (r *FileRequest) EditByID(isTest bool, id int64, file models.File) error {
	return Edit(isTest, "files", r.db, "id", strconv.Itoa(int(id)), file)
}

// GetByID implements Service.
func (r *FileRequest) GetByID(isTest bool, id int64) (models.File, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	file, err := Get[models.File](isTest, "files", r.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *file, err
	}

	return *file, nil
}

// RegisterRoutes implements Service.
func (r *FileRequest) RegisterRoutes() {
	router := r.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	FileRequests := APIV1Router.PathPrefix("/files/").Subrouter()
	FileRequests.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	FileRequests.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	FileRequests.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	FileRequests.HandleFunc("/{id}", middlewares.AdminRoleGuard(r.PatchHandler)).Methods("PATCH")
	FileRequests.HandleFunc("/{id}", middlewares.LoginGuard(r.DeleteHandler)).Methods("DELETE")
	FileRequests.HandleFunc("/{id}/admin", middlewares.AdminRoleGuard(r.DeleteAdminHandler)).Methods("DELETE")
}
