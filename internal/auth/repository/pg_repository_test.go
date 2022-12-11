package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/Yangiboev/auth-example/internal/models"
	"github.com/Yangiboev/auth-example/pkg/utils"
)

func TestAuthRepo_Register(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"name", "password", "email"}).AddRow(
			"dell", "123456", "dell@gmail.com")

		user := &models.User{
			Name:     "icon",
			Email:    "dell@gmail.com",
			Password: "123456",
		}

		mock.ExpectQuery(createUserQuery).WithArgs(&user.Name, &user.Email,
			&user.Password).WillReturnRows(rows)

		createdUser, err := authRepo.Register(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func TestAuthRepo_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetByID", func(t *testing.T) {
		uid := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "name", "email"}).AddRow(
			uid, "dell", "icon", "dell@mail.ru")

		testUser := &models.User{
			UserID: uid,
			Name:   "dell",
			Email:  "dell@mail.ru",
		}

		mock.ExpectQuery(getUserQuery).
			WithArgs(uid).
			WillReturnRows(rows)

		user, err := authRepo.GetByID(context.Background(), uid)
		require.NoError(t, err)
		require.Equal(t, user.Name, testUser.Name)
		fmt.Printf("test user: %s \n", testUser.Name)
		fmt.Printf("user: %s \n", user.Name)
	})
}

func TestAuthRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Delete", func(t *testing.T) {

		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authRepo.Delete(context.Background(), uid)
		require.Nil(t, err)
	})

	t.Run("Delete No rows", func(t *testing.T) {

		uid := uuid.New()

		mock.ExpectExec(deleteUserQuery).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 0))

		err := authRepo.Delete(context.Background(), uid)

		require.NotNil(t, err)
	})
}

func TestAuthRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Update", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "password", "email"}).AddRow(
			"dell", "icon", "123456", "dell@gmail.com")

		user := &models.User{
			Name:     "dell",
			Email:    "dell@gmail.com",
			Password: "123456",
		}

		mock.ExpectQuery(updateUserQuery).WithArgs(&user.Name, &user.Email, &user.UserID).WillReturnRows(rows)

		updatedUser, err := authRepo.Update(context.Background(), user)

		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, user, updatedUser)
	})
}

func TestAuthRepo_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uid := uuid.New()

		rows := sqlmock.NewRows([]string{"user_id", "name", "email"}).AddRow(
			uid, "dell", "icon", "dell@mail.ru")

		testUser := &models.User{
			UserID: uid,
			Name:   "dell",
			Email:  "dell@mail.ru",
		}

		mock.ExpectQuery(findUserByEmail).WithArgs(testUser.Email).WillReturnRows(rows)

		foundUser, err := authRepo.FindByEmail(context.Background(), testUser)

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		require.Equal(t, foundUser.Name, testUser.Name)
	})
}

func TestAuthRepo_GetUsers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByEmail", func(t *testing.T) {
		uid := uuid.New()

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		rows := sqlmock.NewRows([]string{"user_id", "name", "email"}).AddRow(
			uid, "dell", "icon", "dell@mail.ru")

		mock.ExpectQuery(getTotal).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getUsers).WithArgs("", 0, 10).WillReturnRows(rows)

		users, err := authRepo.GetUsers(context.Background(), &utils.PaginationQuery{
			Size:    10,
			Page:    1,
			OrderBy: "",
		})
		require.NoError(t, err)
		require.NotNil(t, users)
	})

}

func TestAuthRepo_FindByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByName", func(t *testing.T) {
		uid := uuid.New()
		userName := "dell"

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		rows := sqlmock.NewRows([]string{"user_id", "name", "email"}).AddRow(
			uid, "dell", "icon", "dell@mail.ru")

		mock.ExpectQuery(getTotalCount).WillReturnRows(totalCountRows)
		mock.ExpectQuery(findUsers).WithArgs("", 0, 10).WillReturnRows(rows)

		usersList, err := authRepo.FindByName(context.Background(), userName, &utils.PaginationQuery{
			Size:    10,
			Page:    1,
			OrderBy: "",
		})

		require.NoError(t, err)
		require.NotNil(t, usersList)
	})
}
