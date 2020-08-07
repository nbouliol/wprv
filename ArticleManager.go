// A simple library to manage articles with Postgres
package ArticleManager

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// ArticleLib describes all the implemented functions
type ArticleLib interface {
	Create(title, body, author string) (Article, error)
	GetAll() ([]Article, error)
	GetById(id int64) (Article, error)
	Update(article Article) (Article, error)
	Delete(id int64) error
}

// Article is the representation of a article
type Article struct {
	Id     int64
	Title  string
	Body   string
	Author string
}

type DB struct {
	*sql.DB
}

// Connect takes a `database/sql` connection string as parameter and returns a db
func Connect(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Create creates a new article into db
func (db *DB) Create(title, body, author string) (Article, error) {
	var insertedId int64 = 0
	sqlStatement := `INSERT INTO articles (title, body, author)
					VALUES ($1, $2, $3)
					RETURNING id`
	err := db.QueryRow(sqlStatement, title, body, author).Scan(&insertedId)

	if err != nil {
		return Article{}, err
	}
	return Article{insertedId, title, body, author}, nil

}

// GetAll returns all the articles into the db
func (db *DB) GetAll() ([]Article, error) {
	rows, err := db.Query("SELECT id, title, body, author FROM articles")
	if err != nil {
		return nil, err
	}
	articles := make([]Article, 0)
	var article Article
	for rows.Next() {
		err := rows.Scan(&article.Id, &article.Title, &article.Body, &article.Author)

		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// GetById returns the article that matches the id
func (db *DB) GetById(id int64) (Article, error) {
	var article Article
	err := db.QueryRow("SELECT id, title, body, author FROM articles WHERE id = $1", id).Scan(&article.Id, &article.Title, &article.Body, &article.Author)

	if err != nil {
		return Article{}, err
	}

	return article, nil
}

// Update updates an article
func (db *DB) Update(article Article) (Article, error) {
	sqlStatement := `UPDATE articles
					SET title = $2, body = $3, author = $4
					WHERE id = $1`
	result, err := db.Exec(sqlStatement, article.Id, article.Title, article.Body, article.Author)

	if err != nil {
		return Article{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Article{}, err
	}
	if count == 0 {
		return Article{}, errors.New("No row updated")
	}

	return article, nil
}

// Delete deletes an article
func (db *DB) Delete(id int64) error {
	_, err := db.Exec("DELETE FROM articles WHERE id = $1", id)

	return err
}
