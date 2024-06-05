package services_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"forum-backend-go/ineternal/models"
	"forum-backend-go/ineternal/services"
	"forum-backend-go/ineternal/utils"
)

func TestUserService(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	database := utils.NewDatabase()
	db, err := database.GetDB(true)
	require.NoError(err)

	userService := services.NewUserService(db, mux.NewRouter())
	want := models.User{Email: "example@example.com", Password: "dddddd", FirstName: "fname", LastName: "lname", Address: "Addffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", PhoneNumber: "09148998933"}
	id, err := userService.Create(true, want)
	require.NoError(err)
	require.NotZero(id)
	get, err := userService.GetByID(true, id)
	require.NoError(err)
	require.Equal(want.Email, get.Email)
	require.Equal(want.FirstName, get.FirstName)
	require.Equal(want.LastName, get.LastName)
	require.Equal(want.Address, get.Address)
	require.Equal(want.PhoneNumber, get.PhoneNumber)

	EditWant := models.User{Email: "example@exle.com", Password: "dddd", FirstName: "ame", LastName: "ame", Address: "Addfffffffffffffffffffffffffffffffffffffffff", PhoneNumber: "09998933"}
	err = userService.EditByID(true, id, EditWant)

	if err != nil {
		require.NoError(err)
	}

	err = userService.DeleteByID(true, id)
	require.NoError(err)
	database.TearDown("users")

	defer db.Close()
}
