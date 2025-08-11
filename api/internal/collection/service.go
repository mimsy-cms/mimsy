package collection

import (
	"context"
)

type Service interface {
	FindById(ctx context.Context, slug string) (*Collection, error)
	FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	FindAll(ctx context.Context) ([]Collection, error)
}

func NewService(collectionRepository Repository) *service {
	return &service{
		collectionRepository: collectionRepository,
	}
}

type service struct {
	collectionRepository Repository
}

func (s *service) FindById(ctx context.Context, slug string) (*Collection, error) {
	return s.collectionRepository.FindBySlug(ctx, slug)
}

func (s *service) FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error) {
	return s.collectionRepository.FindResource(ctx, collection, slug)
}

func (s *service) FindResources(ctx context.Context, collection *Collection) ([]Resource, error) {
	return s.collectionRepository.FindResources(ctx, collection)
}

func (s *service) FindAll(ctx context.Context) ([]Collection, error) {
	return s.collectionRepository.FindAll(ctx)
}
