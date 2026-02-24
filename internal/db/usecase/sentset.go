package usecase

import (
	"time"

	"github.com/dreamervulpi/tourneyBot/internal/db/entity"
)

type SentSetRepo interface {
	Add(setId int64, tournamentPlatform string, messengerPlatform string, tournamentSlug string, sentAt time.Time) (int64, error)
	Get(setId int64) (entity.SentSet, error)
	Del(setId int64) error
	Edit(setId int64, sentAt time.Time) error
}

type SentSet struct {
	Repo SentSetRepo
}

func (s *SentSet) AddSentSet(request entity.SentSetAddRequest) (entity.SentSetAddResponse, error) {
	setId, err := s.Repo.Add(request.SetId, request.TournamentPlatform, request.MessengerPlatform, request.TournamentSlug, request.SentAt)
	if err != nil {
		return entity.SentSetAddResponse{}, err
	}
	return entity.SentSetAddResponse{SetId: setId}, nil
}

func (s *SentSet) EditSentSet(id int64, request entity.SentSetEditRequest) (entity.SentSetEditResponse, error) {
	_, err := s.Repo.Get(id)
	if err != nil {
		return entity.SentSetEditResponse{}, err
	}

	err = s.Repo.Edit(id, request.SentAt)
	if err != nil {
		return entity.SentSetEditResponse{}, err
	}

	return entity.SentSetEditResponse{}, nil
}

func (s *SentSet) DeleteSentSet(id int64) (entity.SentSetDeleteResponse, error) {
	_, err := s.Repo.Get(id)
	if err != nil {
		return entity.SentSetDeleteResponse{}, err
	}

	err = s.Repo.Del(id)
	if err != nil {
		return entity.SentSetDeleteResponse{}, err
	}

	return entity.SentSetDeleteResponse{}, nil
}

func (s *SentSet) GetSentSet(setId int64) (entity.SentSetGetResponse, error) {
	sentSet, err := s.Repo.Get(setId)
	if err != nil {
		return entity.SentSetGetResponse{}, err
	}

	return entity.SentSetGetResponse{
		SetId:              sentSet.SetId,
		TournamentPlatform: sentSet.TournamentPlatform,
		MessengerPlatform:  sentSet.MessengerPlatform,
		TournamentSlug:     sentSet.TournamentSlug,
		SentAt:             sentSet.SentAt}, err
}
