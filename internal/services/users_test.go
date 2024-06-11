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

func TestUserService(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)
	database := utils.NewDatabase()
	db, err := database.GetDB(true)
	require.NoError(err)

	userService := services.NewUserService(db, mux.NewRouter())
	want := models.User{Email: "example@example.com", Password: "dddddd", FirstName: "fname", LastName: "lname", Address: "Addffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", PhoneNumber: "09148998933"}
	id, err := userService.Create(true, want)
	require.NoError(err)
	assert.NotZero(id)
	// Duplicated email
	_, err = userService.Create(true, want)
	require.Error(err)

	get, err := userService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.FirstName, get.FirstName)
	assert.Equal(want.LastName, get.LastName)
	assert.Equal(want.Address, get.Address)
	assert.Equal(want.PhoneNumber, get.PhoneNumber)

	// Email and Password could not change by this commands
	want = models.User{FirstName: "ame", LastName: "ame", Address: "Addfffffffffffffffffffffffffffffffffffffffff", PhoneNumber: "09998933"}
	err = userService.EditByID(true, id, want)

	if err != nil {
		require.NoError(err)
	}

	get, err = userService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.FirstName, get.FirstName)
	assert.Equal(want.LastName, get.LastName)
	assert.Equal(want.Address, get.Address)
	assert.Equal(want.PhoneNumber, get.PhoneNumber)
	// Partially edit
	// Email and Password could not change by this commands
	partialEdited := models.User{Email: "dff"}
	err = userService.EditByID(true, id, partialEdited)

	if err != nil {
		require.NoError(err)
	}

	get, err = userService.GetByID(true, id)
	require.NoError(err)
	assert.Equal(want.FirstName, get.FirstName)
	assert.Equal(want.LastName, get.LastName)
	assert.Equal(want.Address, get.Address)
	assert.Equal(want.PhoneNumber, get.PhoneNumber)

	err = userService.DeleteByID(true, id)
	require.NoError(err)
	database.TearDown("users")

	defer db.Close()
}
