package collection

import (
	"context"
	"encoding/json"
)

type Service interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	FindAll(ctx context.Context, params *FindAllParams) ([]Collection, error)
	FindAllGlobals(ctx context.Context, params *FindAllParams) ([]map[string]any, error)
	DeleteResource(ctx context.Context, resource *Resource) error
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

func (s *service) FindAllGlobals(ctx context.Context, params *FindAllParams) ([]map[string]any, error) {
	globals, err := s.collectionRepository.FindAllGlobals(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	for _, coll := range globals {
		result = append(result, map[string]any{
			"slug":       coll.Slug,
			"name":       coll.Name,
			"fields":     json.RawMessage(coll.Fields),
			"created_at": coll.CreatedAt,
			"created_by": coll.CreatedBy,
			"updated_at": coll.UpdatedAt,
			"updated_by": coll.UpdatedBy,
			"is_global":  coll.IsGlobal,
		})
	}

	return result, nil
}

func (s *service) DeleteResource(ctx context.Context, resource *Resource) error {
	return s.collectionRepository.DeleteResource(ctx, resource)
}
