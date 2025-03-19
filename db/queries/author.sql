-- name: CreateAuthor :exec
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors;

-- name: ListAuthorsByIds :many
SELECT * FROM authors
WHERE id IN (sqlc.slice(ids));

-- name: ListAuthorsPaginated :many
SELECT * FROM authors
LIMIT ?
OFFSET ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;

-- name: DeleteAuthors :exec
DELETE FROM authors
WHERE id IN (sqlc.slice(ids));

-- name: UpdateAuthor :exec
UPDATE authors
SET name = ?, bio = ?
WHERE id = ?;

-- name: GetAuthorWithBooks :many
SELECT 
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    b.id as book_id,
    b.title as book_title
FROM authors a
LEFT JOIN books b ON a.id = b.author_id
WHERE a.id = ?;

