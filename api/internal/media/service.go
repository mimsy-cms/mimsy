package media

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/storage"
)

type MediaService interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader, contentType string, user *auth.User) (*Media, error)
	FindAll(ctx context.Context) ([]Media, error)
}

type mediaService struct {
	storage         storage.Storage
	mediaRepository Repository
}

func NewService(storage storage.Storage, mediaRepository Repository) MediaService {
	return &mediaService{storage: storage, mediaRepository: mediaRepository}
}

func (s *mediaService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, contentType string, user *auth.User) (*Media, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := s.storage.Upload(ctx, id.String(), file, contentType); err != nil {
		return nil, err
	}

	params := &CreateMediaParams{
		Uuid:         id,
		Name:         id.String(),
		ContentType:  contentType,
		Size:         fileHeader.Size,
		UploadedById: user.ID,
	}

	// TODO: Put this in a transaction with the storage upload
	media, err := s.mediaRepository.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *mediaService) FindAll(ctx context.Context) ([]Media, error) {
	return s.mediaRepository.FindAll(ctx)
}
