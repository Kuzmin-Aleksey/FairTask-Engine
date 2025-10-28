package parameters_service

import (
	"FairTask_Engine/internal/domain/entity"
	"context"
	"fmt"
)

type Repo interface {
	Save(ctx context.Context, param *entity.Parameter) error
	GetAll(ctx context.Context) ([]entity.Parameter, error)
	Delete(ctx context.Context, id int) error
}

type ParametersService struct {
	repo Repo
}

func NewParametersService(repo Repo) *ParametersService {
	return &ParametersService{
		repo: repo,
	}
}

func (s *ParametersService) CreateParameter(ctx context.Context, param *entity.Parameter) error {
	const op = "ParametersService.SaveParameter"

	if err := s.repo.Save(ctx, param); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ParametersService) GetAll(ctx context.Context) ([]entity.Parameter, error) {
	const op = "ParametersService.GetAll"
	params, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return params, nil
}

func (s *ParametersService) DeleteParameter(ctx context.Context, id int) error {
	const op = "ParametersService.DeleteParameter"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
