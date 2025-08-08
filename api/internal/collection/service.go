package collection

import (
	"context"
	"encoding/json"
)

func NewService(collectionRepository Repository) *Service {
	return &Service{
		collectionRepository: collectionRepository,
	}
}

type Service struct {
	collectionRepository Repository
}

func (s *Service) GetDefinition(ctx context.Context, slug string) (map[string]any, error) {
	collection, err := s.collectionRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"slug":       slug,
		"name":       collection.Name,
		"fields":     json.RawMessage(collection.Fields),
		"created_at": collection.CreatedAt,
		"created_by": collection.CreatedBy,
		"updated_at": collection.UpdatedAt,
		"updated_by": collection.UpdatedBy,
	}, nil
}

func (s *Service) GetItems(ctx context.Context, slug string) ([]Item, error) {
	exists, err := s.collectionRepository.CollectionExists(ctx, slug)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	return s.collectionRepository.FindItemsBySlug(ctx, slug)
}

func (s *Service) List(ctx context.Context) ([]map[string]any, error) {
	collections, err := s.collectionRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	for _, coll := range collections {
		result = append(result, map[string]any{
			"slug":       coll.Slug,
			"name":       coll.Name,
			"fields":     json.RawMessage(coll.Fields),
			"created_at": coll.CreatedAt,
			"created_by": coll.CreatedBy,
			"updated_at": coll.UpdatedAt,
			"updated_by": coll.UpdatedBy,
		})
	}

	return result, nil
}
