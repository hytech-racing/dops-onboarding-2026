package usecase

import "main/internals/db/repository"

type CarRunUseCase struct {
	repo repository.CarRunRepository
}

func NewCarRunUseCase(repo repository.CarRunRepository) *CarRunUseCase {
	return &CarRunUseCase{
		repo: repo,
	}
}
