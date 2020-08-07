package ArticleManager

import (
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	defer db.Close()

	title := "title"
	body := "body"
	author := "author"
	ddb := &DB{db}

	mock.ExpectQuery("INSERT INTO articles").WithArgs(title, body, author).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))

	// act
	_, err = ddb.Create(title, body, author)
	assert.Nil(t, err)

	// assert
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestCreateWithDuplicateTitle(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	defer db.Close()

	title := "title"
	body := "body"
	author := "author"
	ddb := &DB{db}

	mock.ExpectQuery("INSERT INTO articles").WithArgs(title, body, author).WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	_, err = ddb.Create(title, body, author)
	assert.Nil(t, err)

	// act
	_, err = ddb.Create(title, body, author)
	assert.NotNil(t, err)

	// assert
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)

}

func TestGetAll(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article1, err := ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 2", "body 2", "author 2").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article2, err := ddb.Create("title 2", "body 2", "author 2")
	assert.Nil(t, err)

	rows := sqlmock.NewRows([]string{"id", "title", "body", "author"}).
		AddRow(article1.Id, article1.Title, article1.Body, article1.Author).
		AddRow(article2.Id, article2.Title, article2.Body, article2.Author)

	mock.ExpectQuery("SELECT id, title, body, author FROM articles").WillReturnRows(rows)

	// act
	articles, err := ddb.GetAll()

	// assert
	assert.Nil(t, err)

	assert.Equal(t, 2, len(articles))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetAllEmpty(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "body", "author"})

	mock.ExpectQuery("SELECT id, title, body, author FROM articles").WillReturnRows(rows)

	// act
	articles, err := ddb.GetAll()

	// assert
	assert.Nil(t, err)

	assert.Empty(t, articles)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetById(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	_, err = ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 2", "body 2", "author 2").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article2, err := ddb.Create("title 2", "body 2", "author 2")
	assert.Nil(t, err)

	row := sqlmock.NewRows([]string{"id", "title", "body", "author"}).
		AddRow(article2.Id, article2.Title, article2.Body, article2.Author)

	id := article2.Id

	mock.ExpectQuery("SELECT id, title, body, author FROM articles").WithArgs(id).WillReturnRows(row)

	// act
	article, err := ddb.GetById(id)

	// assert
	assert.Nil(t, err)

	assert.Equal(t, article2, article)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetByIdWithWrongId(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	_, err = ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 2", "body 2", "author 2").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	_, err = ddb.Create("title 2", "body 2", "author 2")
	assert.Nil(t, err)

	row := sqlmock.NewRows([]string{"id", "title", "body", "author"})
	var id int64 = 10

	mock.ExpectQuery("SELECT id, title, body, author FROM articles").WithArgs(id).WillReturnRows(row)

	// act
	article, err := ddb.GetById(id)

	// assert
	assert.NotNil(t, err)

	assert.Empty(t, article)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article, err := ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	newArticle := article
	newArticle.Title = "new title"
	mock.ExpectExec("UPDATE articles").WithArgs(newArticle.Id, newArticle.Title, newArticle.Body, newArticle.Author).WillReturnResult(sqlmock.NewResult(article.Id, 1))

	// act
	updatedArticle, err := ddb.Update(newArticle)
	assert.Nil(t, err)
	assert.NotEqual(t, article, updatedArticle)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateNoRowAffected(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article, err := ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	newArticle := article
	newArticle.Id = 12
	newArticle.Title = "new title"
	mock.ExpectExec("UPDATE articles").WithArgs(newArticle.Id, newArticle.Title, newArticle.Body, newArticle.Author).WillReturnResult(sqlmock.NewResult(article.Id, 0))

	// act
	updatedArticle, err := ddb.Update(newArticle)

	// assert
	assert.EqualError(t, err, "No row updated")
	assert.NotEqual(t, article, updatedArticle)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO articles").WithArgs("title 1", "body 1", "author 1").WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(1))
	article, err := ddb.Create("title 1", "body 1", "author 1")
	assert.Nil(t, err)

	id := article.Id
	mock.ExpectExec("DELETE FROM articles").WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	// act
	err = ddb.Delete(id)
	assert.Nil(t, err)

	// assert
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDeleteWithWrongId(t *testing.T) {
	// arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	ddb := &DB{db}

	defer db.Close()

	var id int64 = -10

	mock.ExpectExec("DELETE FROM articles").WithArgs(id)

	// act
	err = ddb.Delete(id)
	assert.NotNil(t, err)

	// assert
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestConnect(t *testing.T) {
	db, err := Connect("postgres://postgres:p@localhost/articles?sslmode=disable")
	assert.NotNil(t, err)
	assert.Nil(t, db)
}

// Usage example
func Example() {
	db, err := Connect("postgres://postgres:postgres@localhost/articles?sslmode=disable")

	if err != nil {
		log.Panic(err)
	}

	article, err := db.Create("a beatifull title 3", "and a beautifull body", "nico")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(article)
}
