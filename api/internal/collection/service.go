package collection

import (
	"context"
	"encoding/json"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetDefinition(ctx context.Context, slug string) (map[string]interface{}, error) {
	coll, err := s.Repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"slug":       slug,
		"name":       coll.Name,
		"fields":     json.RawMessage(coll.Fields),
		"created_at": coll.CreatedAt,
		"created_by": coll.CreatedBy,
		"updated_at": coll.UpdatedAt,
		"updated_by": coll.UpdatedBy,
	}, nil
}

func (s *Service) GetItems(ctx context.Context, slug string) ([]Item, error) {
	exists, err := s.Repo.CollectionExists(ctx, slug)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}

	return s.Repo.FindItemsBySlug(ctx, slug)
}

func (s *Service) ListCollections(ctx context.Context) ([]map[string]interface{}, error) {
	collections, err := s.Repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, coll := range collections {
		result = append(result, map[string]interface{}{
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
