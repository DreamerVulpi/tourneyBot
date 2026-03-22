package auth

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dreamervulpi/tourneyBot/challonge"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

func GetSessionDiscord() (*discordgo.Session, error) {
	ctx := context.Background()
	dsAuth := &AuthClient{
		Config:    GetDiscordOauth2(),
		TokenFile: "token_discord.json",
	}
	if err := dsAuth.Init(ctx); err != nil {
		return nil, err
	}

	token, err := GetTokenFromFile(dsAuth.TokenFile)
	if err != nil {
		return nil, err
	}
	session, err := discordgo.New("Bot " + token.AccessToken)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func GetClientStartgg() (*startgg.Client, error) {
	ctx := context.Background()
	stAuth := &AuthClient{
		Config:    GetStartggOauth2(),
		TokenFile: "token_startgg.json",
	}
	if err := stAuth.Init(ctx); err != nil {
		return nil, err
	}

	return startgg.NewClient(stAuth.HTTPClient), nil
}
func GetClientChallonge() (*challonge.Client, error) {
	ctx := context.Background()
	chAuth := &AuthClient{
		Config:    GetChallongeOauth2(),
		TokenFile: "token_challonge.json",
	}
	if err := chAuth.Init(ctx); err != nil {
		return nil, err
	}

	token, err := GetTokenFromFile(chAuth.TokenFile)
	if err != nil {
		return nil, err
	}
	return challonge.NewClient(chAuth.HTTPClient, token.AccessToken), nil
}
