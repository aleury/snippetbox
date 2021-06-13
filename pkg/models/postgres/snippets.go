package postgres

import (
	"context"
	"errors"

	"adameury.io/snippetbox/pkg/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SnippetRepo struct {
	DB *pgxpool.Pool
}

func (r *SnippetRepo) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, NOW_UTC(), NOW_UTC() + $3 * '1 DAY'::interval) RETURNING id;`

	row := r.DB.QueryRow(context.Background(), stmt, title, content, expires)

	id := 0
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SnippetRepo) Get(id int) (models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW_UTC() AND id = $1;`

	row := r.DB.QueryRow(context.Background(), stmt, id)

	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return *s, models.ErrNoRecord
		}
		return *s, err
	}

	return *s, nil
}

func (r *SnippetRepo) Latest() ([]models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW_UTC() ORDER BY created DESC LIMIT 10;`

	rows, err := r.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, *s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
