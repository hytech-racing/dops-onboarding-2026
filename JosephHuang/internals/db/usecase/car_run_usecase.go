package usecase

import (
	"context"
	"errors"
	"main/internals/db/repository"
	"main/models"
)

type CarRunUseCase struct {
	repo repository.CarRunRepository
}

func NewCarRunUseCase(repo repository.CarRunRepository) *CarRunUseCase {
	return &CarRunUseCase{
		repo: repo,
	}
}

func (uc *CarRunUseCase) CreateCarRun(ctx context.Context) (*models.CarRun, error) {
	carRun := models.NewCarRun("", "", "", "", models.File{})

	err := uc.repo.Create(ctx, &carRun)
	if err != nil {
		return nil, errors.New("Error: " + err.Error())
	}

	return &carRun, nil
}
