package services

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	searchField := req.URL.Query().Get("search_field")
	searchFieldValue := req.URL.Query().Get("search_field_value")
	files, err := GetMany[models.File](false, "files", r.db, "", "", "", searchField, searchFieldValue, "=", []string{})

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

func (r *File) DownloadHandler(w http.ResponseWriter, req *http.Request) {
	filename := mux.Vars(req)["filename"]
	file, err := os.Open("../../uploads/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.ServeContent(w, req, filename, time.Now(), file)
}

func (r *File) UploadHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	tempString := strings.Split(handler.Filename, ".")
	extension := tempString[len(tempString)-1]
	random := strings.ReplaceAll(uuid.New().String(), "-", "")
	dst, err := os.Create("../../uploads/" + random + "." + extension)
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
		ID       int64  `json:"id"`
		Filename string `json:"filename"`
	}
	user, _ := utils.GetUserFromRequest(req, w)
	NewFile := &models.File{FileName: random + "." + extension, UserID: user.ID}
	id, err := r.Create(false, *NewFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	utils.WriteJSON(w, &data{Filename: random + "." + extension, ID: id})
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

	if user.Role != "admin" && file.UserID != user.ID {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	_ = os.Remove("../../uploads/" + file.FileName)
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
	FilesRouter.HandleFunc("/download/{filename}", r.DownloadHandler).Methods("GET")
	FilesRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	FilesRouter.HandleFunc("/upload", middlewares.LoginGuard(r.UploadHandler)).Methods("POST")
	FilesRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.PatchHandler)).Methods("PATCH")
	FilesRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.DeleteHandler)).Methods("DELETE")
}
