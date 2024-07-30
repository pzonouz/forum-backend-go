package services

import (
	"database/sql"
	"io"
	"log"
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
	operator := req.URL.Query().Get("operator")
	files, err := GetMany[models.File](false, "files", r.db, "", "", "", searchField, searchFieldValue, operator, []string{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, files)
}

func (r *File) GetHandlerForNamed(w http.ResponseWriter, req *http.Request) {
	searchFieldValue := req.URL.Query().Get("title")
	var query string
	if strings.Compare(searchFieldValue, "") == 0 {
		query = `SELECT id,title,filename,created_at,user_id,COALESCE(question_id,0),COALESCE(answer_id,0),filetype FROM files WHERE title is not NULL ORDER BY id DESC LIMIT 50`
	} else {
		query = `SELECT id,title,filename,created_at,user_id,COALESCE(question_id,0),COALESCE(answer_id,0),filetype FROM files WHERE title is not NULL AND title LIKE '%` + searchFieldValue + `%' ORDER BY id DESC  LIMIT 50`
	}
	files := []models.File{}
	var file models.File
	rows, err := r.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&file.ID, &file.Title, &file.FileName, &file.CreatedAt, &file.UserID, &file.QuestionID, &file.AnswerID, &file.FileType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		files = append(files, file)
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

func (r *File) CleanUpHandler(w http.ResponseWriter, req *http.Request) {
	files, err := os.ReadDir("../../uploads/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, file := range files {
		obj, err := Get[models.File](false, "files", r.db, "filename", file.Name(), []string{})
		if err != nil {
			if err.Error() == "not found" {
				err := os.Remove("../../uploads/" + file.Name())
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				goto down
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if obj.QuestionID == 0 && obj.AnswerID == 0 && strings.Compare(obj.Title, "") == 0 {
			err := os.Remove("../../uploads/" + file.Name())
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	down:
	}
	query := `DELETE FROM files WHERE title IS NULL AND question_id IS NULL AND answer_id IS NULL`
	r.db.Exec(query)
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
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
	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		log.Print(err.Error())
		return
	}
	var NewFile *models.File
	NewFile = &models.File{FileName: random + "." + extension, UserID: user.ID}
	var id int64
	questionId, questionErr := strconv.Atoi(req.FormValue("question_id"))
	var query string
	answerId, answerErr := strconv.Atoi(req.FormValue("answer_id"))
	if questionErr != nil && answerErr == nil {
		query = `INSERT INTO files (filename,user_id,answer_id) VALUES($1,$2,$3) RETURNING "id"`
		err = r.db.QueryRow(query, NewFile.FileName, NewFile.UserID, int64(answerId)).Scan(&id)
	}
	if answerErr != nil && questionErr == nil {
		query = `INSERT INTO files (filename,user_id,question_id) VALUES($1,$2,$3) RETURNING "id"`
		err = r.db.QueryRow(query, NewFile.FileName, NewFile.UserID, int64(questionId)).Scan(&id)
	}
	if answerErr != nil && questionErr != nil {
		query = `INSERT INTO files (filename,user_id) VALUES($1,$2) RETURNING "id"`
		err = r.db.QueryRow(query, NewFile.FileName, NewFile.UserID).Scan(&id)
	}
	if err != nil {
		print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, &data{Filename: random + "." + extension, ID: id})
}

func (r *File) UploadAdminHandler(w http.ResponseWriter, req *http.Request) {
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
		Title    string `json:"title"`
	}
	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		log.Print(err.Error())
		return
	}
	title := req.FormValue("title")
	var NewFile *models.File
	NewFile = &models.File{FileName: random + "." + extension, UserID: user.ID, Title: title}
	var id int64
	query := `INSERT INTO files (filename,user_id,title) VALUES($1,$2,$3) RETURNING "id"`
	err = r.db.QueryRow(query, NewFile.FileName, NewFile.UserID, NewFile.Title).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, &data{Filename: random + "." + extension, ID: id})
}

func (r *File) IsUniqueFileNameHandler(w http.ResponseWriter, req *http.Request) {
	type File struct {
		Title string `json:"title"`
		ID    int64  `json:"id"`
	}
	file := utils.ReadJSON[File](w, req)
	var count int
	query := `SELECT COUNT(*) FROM files WHERE title=$1 AND id!=$2`
	err := r.db.QueryRow(query, file.Title, file.ID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if count > 0 {
		http.Error(w, "", http.StatusConflict)
		return
	}
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

func (r *File) DeleteAdminHandler(w http.ResponseWriter, req *http.Request) {
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
	FilesRouter.HandleFunc("/collection/", r.GetHandlerForNamed).Methods("GET")
	FilesRouter.HandleFunc("/clean_up", middlewares.AdminRoleGuard(r.CleanUpHandler)).Methods("GET")
	FilesRouter.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	FilesRouter.HandleFunc("/download/{filename}", r.DownloadHandler).Methods("GET")
	FilesRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	FilesRouter.HandleFunc("/upload", middlewares.LoginGuard(r.UploadHandler)).Methods("POST")
	FilesRouter.HandleFunc("/upload_admin", middlewares.AdminRoleGuard(r.UploadAdminHandler)).Methods("POST")
	FilesRouter.HandleFunc("/is_filename_unique", middlewares.AdminRoleGuard(r.IsUniqueFileNameHandler)).Methods("POST")
	FilesRouter.HandleFunc("/{id}", middlewares.AdminRoleGuard(r.PatchHandler)).Methods("PATCH")
	FilesRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.DeleteHandler)).Methods("DELETE")
	FilesRouter.HandleFunc("/{id}/admin", middlewares.AdminRoleGuard(r.DeleteAdminHandler)).Methods("DELETE")
}
