package usecase

import (
	"time"

	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
)

type ParticipantRepo interface {
	Add(gamerTag string,
		messengerPlatform string,
		messengerPlatformId string,
		messenagerPlatformLogin string,
		updatedAt time.Time,
		isFound bool,
		locale string) (string, string, error)
	Edit(gamerTag string,
		messengerPlatform string,
		messengerPlatformId string,
		messengerPlatformLogin string,
		updatedAt time.Time,
		isFound bool,
		locale string) error
	Del(gamerTag string,
		messengerPlatform string) error
	Get(gamerTag string,
		messengerPlatform string) (entity.Participant, error)
}

type Participant struct {
	Repo ParticipantRepo
}

func (p *Participant) AddParticipant(request entity.ParticipantAddRequest) (entity.ParticipantAddResponse, error) {
	gamerTag, messenagerPlatform, err := p.Repo.Add(
		request.GamerTag,
		request.MessengerPlatform,
		request.MessengerPlatformId,
		request.MessengerPlatformLogin,
		request.UpdatedAt,
		request.IsFound,
		request.Locale,
	)
	if err != nil {
		return entity.ParticipantAddResponse{}, err
	}
	return entity.ParticipantAddResponse{GamerTag: gamerTag, MessengerPlatform: messenagerPlatform}, nil
}

func (p *Participant) EditParticipant(request entity.ParticipantEditRequest) (entity.ParticipantEditResponse, error) {
	_, err := p.Repo.Get(request.GamerTag, request.MessenagerPlatform)
	if err != nil {
		return entity.ParticipantEditResponse{}, err
	}

	err = p.Repo.Edit(
		request.GamerTag,
		request.MessenagerPlatform,
		request.MessenagerPlatformId,
		request.MessenagerPlatformLogin,
		request.UpdatedAt,
		request.IsFound,
		request.Locale,
	)
	if err != nil {
		return entity.ParticipantEditResponse{}, err
	}
	return entity.ParticipantEditResponse{}, nil
}

func (p *Participant) DelParticipant(request entity.ParticipantDeleteRequest) (entity.ParticipantDeleteResponse, error) {
	_, err := p.Repo.Get(request.GamerTag, request.MessenagerPlatform)
	if err != nil {
		return entity.ParticipantDeleteResponse{}, err
	}

	err = p.Repo.Del(request.GamerTag, request.MessenagerPlatform)
	if err != nil {
		return entity.ParticipantDeleteResponse{}, err
	}
	return entity.ParticipantDeleteResponse{}, nil
}

func (p *Participant) GetParticipant(request entity.ParticipantGetRequest) (entity.ParticipantGetResponse, error) {
	participant, err := p.Repo.Get(request.GamerTag, request.MessenagerPlatform)
	if err != nil {
		return entity.ParticipantGetResponse{}, err
	}
	return entity.ParticipantGetResponse{
		GamerTag:               participant.GamerTag,
		MessengerPlatform:      participant.MessengerPlatform,
		MessengerPlatformId:    participant.MessengerPlatformId,
		MessengerPlatformLogin: participant.MessengerPlatformLogin,
		UpdatedAt:              participant.UpdatedAt,
	}, nil
}
