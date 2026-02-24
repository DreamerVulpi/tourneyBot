-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS participants (
    gamer_tag TEXT NOT NULL,
    platform TEXT NOT NULL,
    platform_id TEXT NOT NULL,
    platform_login TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (gamer_tag, platform)
);

CREATE TABLE IF NOT EXISTS sent_sets (
    set_id BIGINT NOT NULL,
    tournament_platform TEXT NOT NULL,
    messenger_platform TEXT NOT NULL,
    tournament_slug TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(set_id, tournament_platform, messenger_platform)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sent_sets;
DROP TABLE IF EXISTS participants;
-- +goose StatementEnd
