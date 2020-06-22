package repository

import (
	"database/sql"
	"fmt"
	"regexp"
	"summer-web/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	userRepo UserRepository
	mock     sqlmock.Sqlmock
	db       *sql.DB
	gdb      *gorm.DB
	err      error
)

func setup() {
	db, mock, err = sqlmock.New()

	if err != nil {
		fmt.Println(err.Error())
	}

	gdb, err = gorm.Open("postgres", db)

	if err != nil {
		fmt.Println(err.Error())
	}

	userRepo = NewUserRepository(gdb)
}

func TestGetUserByID(t *testing.T) {
	setup()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "password", "follower_count", "following_count", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "test1", "test1", "test1@test1.com", "test1", 1, 1, time.Now(), time.Now(), nil)

	const sqlSelectByID = `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = $1))`
	const searchID = uint(1)

	user := models.User{}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).WithArgs(searchID).WillReturnRows(rows)

	err := userRepo.GetUserByID(searchID, &user)

	assert.Nil(t, err)
	assert.Equal(t, searchID, user.ID)
}

func TestGetUserByUsername(t *testing.T) {
	setup()

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "password", "follower_count", "following_count", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "test1", "test1", "test1@test1.com", "test1", 1, 1, time.Now(), time.Now(), nil)

	const sqlSelectByUsername = `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((username = $1))`
	const username = "test1"

	user := models.User{}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByUsername)).WillReturnRows(rows)

	err := userRepo.GetUserByUsername(username, &user)

	assert.Nil(t, err)
	assert.Equal(t, username, user.Username)
}

func TestAddUser(t *testing.T) {
	setup()

	user := models.User{Username: "test1", Name: "test1", Email: "test1@test1.com", Password: "test1", FollowerCount: 1, FollowingCount: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: nil}
	newID := uint(1)
	const sqlInsert = `INSERT INTO "users" ("username","name","email","password","follower_count","following_count","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "users"."id"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).WithArgs(user.Username, user.Name, user.Email, user.Password, user.FollowerCount, user.FollowingCount, user.CreatedAt, user.UpdatedAt, user.DeletedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()

	assert.Equal(t, uint(0), user.ID)

	err := userRepo.AddUser(&user)

	assert.Nil(t, err)
	assert.Equal(t, newID, user.ID)
}

func TestUpdateUser(t *testing.T) {
	setup()

	user := models.User{ID: 1, Username: "changed"}
	const sqlUpdate = `UPDATE "users" SET "id" = $1, "updated_at" = $2, "username" = $3 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $4`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).WithArgs(user.ID, time.Now(), user.Username, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.UpdateUser(user)

	assert.Nil(t, err)
}
