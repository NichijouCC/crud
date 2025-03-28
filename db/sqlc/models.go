// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"
)

type Author struct {
	ID   int64          `db:"id" json:"id"`
	Name string         `db:"name" json:"name"`
	Bio  sql.NullString `db:"bio" json:"bio"`
}

type Book struct {
	ID       int64  `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	AuthorID int64  `db:"author_id" json:"author_id"`
}
