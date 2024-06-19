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

func TestRoleService(t *testing.T) {
	// t.Parallel()
	assert := assert.New(t)
	require := require.New(t)
	database := utils.NewDatabase()
	db, err := database.GetDB(true)
	require.NoError(err)

	roleService := services.NewRoleService(db, mux.NewRouter())
	want := models.Role{Name: "NewRole"}
	id, err := roleService.Create(true, want)
	require.NoError(err)
	assert.NotZero(id)

	_, err = roleService.Create(true, want)
	require.Error(err)

	get, err := roleService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Name, get.Name)

	want = models.Role{Name: "ame"}
	err = roleService.EditByID(true, id, want)

	if err != nil {
		require.NoError(err)
	}

	get, err = roleService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Name, get.Name)
	partialEdited := models.Role{Name: want.Name}
	err = roleService.EditByID(true, id, partialEdited)

	if err != nil {
		require.NoError(err)
	}

	get, err = roleService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.Name, get.Name)

	err = roleService.DeleteByID(true, id)
	require.NoError(err)

	defer db.Close()
}
