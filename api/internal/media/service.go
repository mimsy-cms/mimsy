package media

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/mimsy-cms/mimsy/internal/storage"
)

type MediaService interface {
	Upload(ctx context.Context, file multipart.File, contentType string) error
}

type mediaService struct {
	storage storage.Storage
}

func NewMediaService(storage storage.Storage) MediaService {
	return &mediaService{storage: storage}
}

func (s *mediaService) Upload(ctx context.Context, file multipart.File, contentType string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	if err := s.storage.Upload(ctx, id.String(), file, contentType); err != nil {
		return err
	}

	return nil
}
