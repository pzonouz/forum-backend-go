package services

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"forum-backend-go/internal/middlewares"
	"forum-backend-go/internal/models"
	"forum-backend-go/internal/utils"
)

func NewFileService(db *sql.DB, router *mux.Router) *File {
	return &File{
		db:     db,
		router: router,
	}
}

type File struct {
	db     *sql.DB
	router *mux.Router
}

func (r *File) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
	files, err := GetMany[models.File](false, "files", r.db, "", "", "", "", "", "", []string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, files)
}

func (r *File) GetHandler(w http.ResponseWriter, req *http.Request) {
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

func (r *File) UploadHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	random := strings.ReplaceAll(uuid.New().String(), "-", "")
	dst, err := os.Create("../../uploads/" + random)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type data struct {
		Location string `json:"location"`
	}
	utils.WriteJSON(w, &data{Location: "uploads/" + random})

	// id, err := r.Create(false, file)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }
	//
	// type data struct {
	// 	ID int64 `json:"id"`
	// }
	//
	// utils.WriteJSON(w, &data{ID: id})
}

func (r *File) PostHandler(w http.ResponseWriter, req *http.Request) {
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

func (r *File) PatchHandler(w http.ResponseWriter, req *http.Request) {
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

	if user.Role != "admin" {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if filePartial.Name != "" && len(filePartial.Name) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	err = r.EditByID(false, int64(id), filePartial)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *File) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}

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

	if user.Role != "admin" {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	err = r.DeleteByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// Create implements Service.
func (r *File) Create(isTest bool, file models.File) (int64, error) {
	var id int64
	id, err := Create[models.File](isTest, "files", file, r.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (r *File) DeleteByID(isTest bool, id int64) error {
	return Delete[models.File](isTest, "files", r.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (r *File) EditByID(isTest bool, id int64, file models.File) error {
	return Edit(isTest, "files", r.db, "id", strconv.Itoa(int(id)), file)
}

// GetByID implements Service.
func (r *File) GetByID(isTest bool, id int64) (models.File, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	file, err := Get[models.File](isTest, "files", r.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *file, err
	}

	return *file, nil
}

// RegisterRoutes implements Service.
func (r *File) RegisterRoutes() {
	router := r.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	FilesRouter := APIV1Router.PathPrefix("/files/").Subrouter()
	FilesRouter.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	FilesRouter.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	FilesRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	FilesRouter.HandleFunc("/upload", middlewares.LoginGuard(r.UploadHandler)).Methods("POST")
	FilesRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.PatchHandler)).Methods("PATCH")
	FilesRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.DeleteHandler)).Methods("DELETE")
}
