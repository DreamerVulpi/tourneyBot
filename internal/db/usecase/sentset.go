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
	Exists(setId int64) (bool, error)
}

type SentSet struct {
	Repo SentSetRepo
}

func (s *SentSet) IsExists(request entity.SentSetCheckRequest) (entity.SentSetCheckResponse, error) {
	state, err := s.Repo.Exists(request.SetId)
	if err != nil {
		return entity.SentSetCheckResponse{}, err
	}
	return entity.SentSetCheckResponse{State: state}, nil
}

func (s *SentSet) AddSentSet(request entity.SentSetAddRequest) (entity.SentSetAddResponse, error) {
	setId, err := s.Repo.Add(request.SetId, request.TournamentPlatform, request.MessengerPlatform, request.TournamentSlug, request.SentAt)
	if err != nil {
		return entity.SentSetAddResponse{}, err
	}
	return entity.SentSetAddResponse{SetId: setId}, nil
}

func (s *SentSet) EditSentSet(request entity.SentSetEditRequest) (entity.SentSetEditResponse, error) {
	_, err := s.Repo.Get(request.SetId)
	if err != nil {
		return entity.SentSetEditResponse{}, err
	}

	err = s.Repo.Edit(request.SetId, request.SentAt)
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
