package collection

import (
	"context"
)

type Service interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	FindAll(ctx context.Context, params *FindAllParams) ([]Collection, error)
	CreateResource(ctx context.Context, collection *Collection, resourceSlug string, createdBy int64) (*Resource, error)
	FindAllGlobals(ctx context.Context, params *FindAllParams) ([]Collection, error)
	UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error)
	DeleteResource(ctx context.Context, resource *Resource) error
	FindUserEmail(ctx context.Context, id int64) (string, error)
}

func NewService(collectionRepository Repository) *service {
	return &service{
		collectionRepository: collectionRepository,
	}
}

type service struct {
	collectionRepository Repository
}

func (s *service) FindBySlug(ctx context.Context, slug string) (*Collection, error) {
	return s.collectionRepository.FindBySlug(ctx, slug)
}

func (s *service) FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error) {
	return s.collectionRepository.FindResource(ctx, collection, slug)
}

func (s *service) FindResources(ctx context.Context, collection *Collection) ([]Resource, error) {
	return s.collectionRepository.FindResources(ctx, collection)
}

func (s *service) FindAll(ctx context.Context, params *FindAllParams) ([]Collection, error) {
	return s.collectionRepository.FindAll(ctx, params)
}

func (s *service) CreateResource(ctx context.Context, collection *Collection, resourceSlug string, createdBy int64) (*Resource, error) {
	return s.collectionRepository.CreateResource(ctx, collection, resourceSlug, createdBy)
}

func (s *service) FindAllGlobals(ctx context.Context, params *FindAllParams) ([]Collection, error) {
	return s.collectionRepository.FindAllGlobals(ctx, params)
}

func (s *service) UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error) {
	return s.collectionRepository.UpdateResourceContent(ctx, collection, resourceSlug, content)
}

func (s *service) DeleteResource(ctx context.Context, resource *Resource) error {
	return s.collectionRepository.DeleteResource(ctx, resource)
}

func (s *service) FindUserEmail(ctx context.Context, id int64) (string, error) {
	return s.collectionRepository.FindUserEmail(ctx, id)
}
