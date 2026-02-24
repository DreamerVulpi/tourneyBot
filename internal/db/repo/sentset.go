package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SentSet struct {
	Conn *pgxpool.Pool
}

func (s *SentSet) Exists(setId int64) (bool, error) {
	const sql = `
		SELECT EXISTS
		(SELECT 1 FROM sent_sets WHERE set_id = $1)
	`
	var exists bool
	err := s.Conn.QueryRow(context.Background(), sql, setId).Scan(&exists)
	return exists, err
}

func (s *SentSet) Add(setId int64, tournamentPlatform string, messengerPlatform string, tournamentSlug string, sentAt time.Time) (int, error) {
	const sql = `
		INSERT INTO sent_sets
			(set_id, tournament_platform, messenger_platform, tournament_slug, sent_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (set_id, tournament_platform, messenger_platform) DO NOTHING
		RETURNING set_id`
	var id int
	err := s.Conn.QueryRow(context.Background(), sql, setId, tournamentPlatform, messengerPlatform, tournamentSlug, sentAt).Scan(&id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil // todo: уже существует
		}
		return 0, fmt.Errorf("unable to create sentSet in database, %w", err)
	}
	return id, nil
}

func (s *SentSet) Get(setId int64) (entity.SentSet, error) {
	const sql = `
		SELECT s.set_id, s.tournament_platform, s.messenger_platform, s.tournament_slug, s.sent_at
		FROM sent_sets s
		WHERE set_id = $1`
	var sentSet entity.SentSet
	err := s.Conn.QueryRow(context.Background(), sql, setId).Scan(
		&sentSet.SetId,
		&sentSet.TournamentPlatform,
		&sentSet.MessengerPlatform,
		&sentSet.TournamentSlug,
		&sentSet.SentAt,
	)
	if err != nil {
		return entity.SentSet{}, fmt.Errorf("unable to get sentSet in database, %w", err)
	}
	return sentSet, nil
}

func (s *SentSet) Del(setId int64) error {
	const sql = `
		DELETE FROM sent_sets
		WHERE set_id = $1`
	tag, err := s.Conn.Exec(context.Background(), sql, setId)
	if err != nil {
		return fmt.Errorf("don't deleted sentset from database, %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sentset doesn't exist")
	}

	return nil
}

func (s *SentSet) Edit(id int64, sentAt time.Time) error {
	const sql = `
		UPDATE sent_sets
		SET sent_at = $1
		WHERE set_id = $2`
	tag, err := s.Conn.Exec(context.Background(), sql, sentAt, id)
	if err != nil {
		return fmt.Errorf("don't edited sentset from database, %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sentset doesn't exist")
	}

	return nil
}
