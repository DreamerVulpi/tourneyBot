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
		updatedAt time.Time) (string, string, error)
	Edit(id int, gamerTag string,
		messengerPlatform string,
		messengerPlatformId string,
		messengerPlatformLogin string,
		updatedAt time.Time) error
	Del(id int) error
	Get(id int) (entity.Participant, error)
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
		request.UpdatedAt)
	if err != nil {
		return entity.ParticipantAddResponse{}, err
	}
	return entity.ParticipantAddResponse{GamerTag: gamerTag, MessengerPlatform: messenagerPlatform}, nil
}

func (p *Participant) EditParticipant(id int, request entity.ParticipantEditRequest) (entity.ParticipantEditResponse, error) {
	_, err := p.Repo.Get(id)
	if err != nil {
		return entity.ParticipantEditResponse{}, err
	}

	err = p.Repo.Edit(
		id,
		request.GamerTag,
		request.MessenagerPlatform,
		request.MessenagerPlatformId,
		request.MessenagerPlatformLogin,
		request.UpdatedAt)
	if err != nil {
		return entity.ParticipantEditResponse{}, err
	}
	return entity.ParticipantEditResponse{}, nil
}

func (p *Participant) DelParticipant(id int) (entity.ParticipantDeleteResponse, error) {
	_, err := p.Repo.Get(id)
	if err != nil {
		return entity.ParticipantDeleteResponse{}, err
	}

	err = p.Repo.Del(id)
	if err != nil {
		return entity.ParticipantDeleteResponse{}, err
	}
	return entity.ParticipantDeleteResponse{}, nil
}

func (p *Participant) GetParticipant(id int) (entity.ParticipantGetResponse, error) {
	participant, err := p.Repo.Get(id)
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
