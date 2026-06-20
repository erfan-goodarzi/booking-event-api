package store

import (
	"database/sql"

	"github.com/erfan-goodarzi/booking-event-api/internal/db"
	"github.com/erfan-goodarzi/booking-event-api/internal/models"
)

type PlaylistStore interface {
	GetAll() ([]models.Playlist, error)
	GetById(id string) (*models.Playlist, error)
	Create(p *models.Playlist) (*models.Playlist, error)
	Update(*models.Playlist) (*models.Playlist, error)
	Delete(id string) error
	GetOwner(playlistId string) (string, error)
}

type PostgresPlaylistStore struct {
	db db.DB
}

func NewPostgresPlaylistStore(db db.DB) *PostgresPlaylistStore {
	return &PostgresPlaylistStore{db: db}
}

func (pg *PostgresPlaylistStore) GetAll() ([]models.Playlist, error) {
	var playlists []models.Playlist

	query := `
	SELECT id, user_id, name, color, created_at, updated_at
	FROM playlists
	`

	rows, err := pg.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Playlist

		err := rows.Scan(
			&p.ID,
			&p.UserId,
			&p.Name,
			&p.Color,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		playlists = append(playlists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (pg *PostgresPlaylistStore) GetById(id string) (*models.Playlist, error) {
	var p models.Playlist

	err := pg.db.QueryRow(`
    SELECT id, user_id, name, color, created_at, updated_at
    FROM playlists
    WHERE id = $1
`, id).Scan(
		&p.ID,
		&p.UserId,
		&p.Name,
		&p.Color,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	query := `
  SELECT e.id,
    e.title,
    e.description,
    e.location,
    e.date_time,
    e.user_id,
    e.duration,
		e.version,
    e.created_at,
    e.updated_at
 	FROM events e
	JOIN playlist_events pe
  ON pe.event_id = e.id
	WHERE pe.playlist_id = $1;
 `
	rows, err := pg.db.Query(query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Event

		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Location,
			&e.DateTime,
			&e.UserId,
			&e.Duration,
			&e.Version,
			&e.CreatedAt,
			&e.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		p.Events = append(p.Events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (pg *PostgresPlaylistStore) Create(p *models.Playlist) (*models.Playlist, error) {
	query := `
	INSERT INTO playlists(user_id, name, color)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
  `

	err := pg.db.QueryRow(query, p.UserId, p.Name, p.Color).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (pg *PostgresPlaylistStore) Update(p *models.Playlist) (*models.Playlist, error) {
	query := `UPDATE playlists
	SET name = $1, color = $2
	WHERE id = $3
	RETURNING updated_at
	`

	err := pg.db.QueryRow(query, p.Name, p.Color, p.ID).Scan(&p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return p, nil
}

func ApplyPlaylistPatch(p *models.Playlist, patch models.PatchPlaylistRequest) error {
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Color != nil {
		p.Color = *patch.Color
	}
	return nil
}

func (pg *PostgresPlaylistStore) Delete(id string) error {
	query := `DELETE FROM playlists WHERE id = $1`

	stmt, err := pg.db.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresPlaylistStore) GetOwner(playlistId string) (string, error) {
	var userID string

	query := `
  SELECT user_id
  FROM playlists
  WHERE id = $1
  `

	err := pg.db.QueryRow(query, playlistId).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
