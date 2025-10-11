package usecase

import (
	"context"
	"main/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockCarRunRepository struct {
	mock.Mock
}

func (m *MockCarRunRepository) Create(ctx context.Context, carRun *models.CarRun) error {
	args := m.Called(ctx, carRun)
	return args.Error(0)
}

func (m *MockCarRunRepository) Update(ctx context.Context, carRun *models.CarRun, id primitive.ObjectID) error {
	args := m.Called(ctx, carRun, id)
	return args.Error(0)
}

func (m *MockCarRunRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCarRunRepository) List(ctx context.Context) ([]*models.CarRun, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.CarRun), args.Error(1)
}

func TestCreateCarRun(t *testing.T) {
	mockRepo := new(MockCarRunRepository)
	useCase := NewCarRunUseCase(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CarRun")).Return(nil)

	carRun, err := useCase.CreateCarRun(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, carRun)
	mockRepo.AssertExpectations(t)
}
