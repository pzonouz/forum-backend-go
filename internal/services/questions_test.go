package services_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"forum-backend-go/internal/models"
	"forum-backend-go/internal/services"
	"forum-backend-go/internal/utils"
)

func TestQuestionService(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)
	require := require.New(t)
	database := utils.NewDatabase()
	db, err := database.GetDB(true)
	require.NoError(err)

	questionService := services.NewQuestionService(db, mux.NewRouter())
	want := models.Question{Title: "newTile", Description: "NewDescription"}
	id, err := questionService.Create(true, want)
	require.NoError(err)
	assert.NotZero(id)
	// Duplicated email
	_, err = questionService.Create(true, want)
	require.Error(err)

	get, err := questionService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Title, get.Title)
	assert.Equal(want.Description, get.Description)

	want = models.Question{Title: "editedTitle", Description: "EditedDescriptoim"}
	err = questionService.EditByID(true, id, want)

	if err != nil {
		require.NoError(err)
	}

	get, err = questionService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Title, get.Title)
	assert.Equal(want.Description, get.Description)
	// Partially edit
	// Email and Password could not change by this commands
	partialEdited := models.Question{Title: want.Title}
	err = questionService.EditByID(true, id, partialEdited)

	if err != nil {
		require.NoError(err)
	}

	get, err = questionService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Title, get.Title)
	assert.Equal(want.Description, get.Description)

	err = questionService.DeleteByID(true, id)
	require.NoError(err)
	database.TearDown("questions")

	defer db.Close()
}
