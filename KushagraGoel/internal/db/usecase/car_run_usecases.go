package usecase

import (
	"KushagraGoel/internal/db/repository"
	"KushagraGoel/models"
	"context"
	"fmt"
)


type CarRunUseCase struct {
	carRunRepo repository.CarRunRepository
}

func NewCarRunUseCase(carRunRepo repository.CarRunRepository) *CarRunUseCase {
	return &CarRunUseCase{
		carRunRepo: carRunRepo,
	}
}

func (uc *CarRunUseCase) Create(ctx context.Context, location, car_model, event_type, notes, aws_bucket, file_path, file_name string) (*models.Car_Run, error){
	carRun, err := models.NewCar(location, car_model, event_type, notes, aws_bucket, file_path, file_name)
	if err != nil {
		return nil, fmt.Errorf("failed ot create car model: %w", err)
	}

	if err :=uc.carRunRepo.Create(ctx, carRun); err != nil {
		return nil, fmt.Errorf("failed t persist car run: %w", err)
	}

	return carRun, nil
}