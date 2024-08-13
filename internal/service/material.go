package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MXkodo/Management-of-School-materials/internal/repo"
	"github.com/MXkodo/Management-of-School-materials/model"
	"github.com/google/uuid"
)

type MaterialService interface {
	CreateMaterial(ctx context.Context, mat *model.Material) (string, error)
	GetMaterial(ctx context.Context, uuid string) (*model.Material, error)
	UpdateMaterial(ctx context.Context, mat *model.Material) error
	GetAllMaterials(ctx context.Context, filter model.MaterialFilter, limit, offset int) ([]model.Material, error)
}

type MaterialServiceImpl struct {
	repo repo.MaterialRepository
}

func NewMaterialService(repo repo.MaterialRepository) *MaterialServiceImpl {
	return &MaterialServiceImpl{repo: repo}
}

func (s *MaterialServiceImpl) CreateMaterial(ctx context.Context, mat *model.Material) (string, error) {
	mat.UUID = uuid.New().String()
	mat.CreatedAt = time.Now()
	mat.UpdatedAt = time.Now()

	return s.repo.CreateMaterial(ctx, mat)
}

func (s *MaterialServiceImpl) GetMaterial(ctx context.Context, uuid string) (*model.Material, error) {
	return s.repo.GetMaterial(ctx, uuid)
}

func (s *MaterialServiceImpl) UpdateMaterial(ctx context.Context, mat *model.Material) error {
	existingMaterial, err := s.repo.GetMaterial(ctx, mat.UUID)
	if err != nil {
		return err
	}

	if existingMaterial.Type != mat.Type {
		return fmt.Errorf("cannot change material type")
	}

	mat.UpdatedAt = time.Now()
	return s.repo.UpdateMaterial(ctx, mat)
}

func (s *MaterialServiceImpl) GetAllMaterials(ctx context.Context, filter model.MaterialFilter, limit, offset int) ([]model.Material, error) {
	return s.repo.GetAllMaterials(ctx, filter, limit, offset)
}
