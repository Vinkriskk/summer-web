package repository

import (
	"database/sql"
	"fmt"
	"regexp"
	"summer-web/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var (
	postRepo PostRepository
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

	postRepo = NewPostRepository(gdb)
}

func TestGetPosts(t *testing.T) {
	setup()

	rows := sqlmock.NewRows([]string{"id", "caption", "user_id"}).AddRow(1, "hello1", 1).AddRow(2, "hello2", 2)

	const sqlSelectAll = `SELECT * FROM "posts"`

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).WillReturnRows(rows)

	posts, err := postRepo.GetPosts()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(posts))
}

func TestAddPost(t *testing.T) {
	setup()

	post := models.Post{Caption: "123", UserID: 1}
	const sqlInsert = `INSERT INTO "posts" ("caption","user_id") VALUES ($1,$2) RETURNING "posts"."id"`
	newID := uint(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).WithArgs(post.Caption, post.UserID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newID))
	mock.ExpectCommit()

	assert.Equal(t, uint(0), post.ID)

	err := postRepo.AddPost(&post)

	assert.Nil(t, err)
	assert.Equal(t, newID, post.ID)
}
