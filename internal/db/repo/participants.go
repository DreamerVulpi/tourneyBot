package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Participants struct {
	Conn *pgxpool.Pool
}

func (p *Participants) Add(gamerTag string, messenagerPlatform string, messenagerPlatformId string, messenagerPlatformLogin string, updatedAt time.Time) (string, string, error) {
	const sql = `
		INSERT INTO participants
			(gamer_tag, platform, platform_id, platform_login, updated_at)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING gamer_tag, platform`
	var resGamerTag, resMessenagerPlatform string
	err := p.Conn.QueryRow(context.Background(), sql, gamerTag, messenagerPlatform, messenagerPlatformId, messenagerPlatformLogin, updatedAt).Scan(&resGamerTag, &resMessenagerPlatform)
	if err != nil {
		return "", "", fmt.Errorf("unable to create participants in database, %w", err)
	}
	return resGamerTag, resMessenagerPlatform, nil
}

func (p *Participants) Edit(gamerTag string, messenagerPlatform string, messenagerPlatformId string, messenagerPlatformLogin string, updatedAt time.Time) error {
	const sql = `
		UPDATE participants
		SET platform_id = $3, platform_login = $4, updated_at = $5
		WHERE gamer_tag = $1 AND platform = $2`

	tag, err := p.Conn.Exec(context.Background(), sql, gamerTag, messenagerPlatform, messenagerPlatformId, messenagerPlatformLogin, updatedAt)
	if err != nil {
		return fmt.Errorf("don't edited participant from database, %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("participant doesn't exist")
	}

	return nil
}

func (p *Participants) Get(gamerTag, messenagerPlatform string) (entity.Participant, error) {
	const sql = `
		SELECT p.gamer_tag, p.platform, p.platform_id, p.platform_login, p.updated_at
		FROM participants p
		WHERE gamer_tag = $1 AND platform = $2`

	var participant entity.Participant
	err := p.Conn.QueryRow(context.Background(), sql, gamerTag, messenagerPlatform).Scan(
		&participant.GamerTag,
		&participant.MessengerPlatform,
		&participant.MessengerPlatformId,
		&participant.MessengerPlatformLogin,
		&participant.UpdatedAt)
	if err != nil {
		return entity.Participant{}, fmt.Errorf("unable to create participant in database, %w", err)
	}
	return participant, nil
}

func (p *Participants) Del(gamerTag, platform string) error {
	const sql = `
		DELETE FROM participants
		WHERE gamer_tag = $1 AND platform = $2`
	tag, err := p.Conn.Exec(context.Background(), sql, gamerTag, platform)
	if err != nil {
		return fmt.Errorf("don't deleted participant from database, %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("participant doesn't exist")
	}

	return nil
}
