package usecase

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"
)

type AwardFormLogUseCase interface {
	CreateLog(ctx context.Context, changedBy uint, req *awardformdto.CreateAwardFormLogRequest) (*models.AwardFormLog, error)
	GetLogsByFormID(ctx context.Context, formID uint) ([]models.AwardFormLog, error)
}

type awardFormLogUseCase struct {
	repo *repository.AwardFormLogRepository
}

func NewAwardFormLogUseCase(repo *repository.AwardFormLogRepository) AwardFormLogUseCase {
	return &awardFormLogUseCase{repo: repo}
}

func (u *awardFormLogUseCase) CreateLog(ctx context.Context, changedBy uint, req *awardformdto.CreateAwardFormLogRequest) (*models.AwardFormLog, error) {
	userID := int(changedBy)
	now := time.Now()
	log := &models.AwardFormLog{
		FormID:     req.FormID,
		FieldName:  req.FieldName,
		OldValue:   req.OldValue,
		NewValue:   req.NewValue,
		ChangedBy:  userID,
		CreatedAt:  now,
	}

	if err := u.repo.Create(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func (u *awardFormLogUseCase) GetLogsByFormID(ctx context.Context, formID uint) ([]models.AwardFormLog, error) {
	return u.repo.GetByFormID(ctx, formID)
}
