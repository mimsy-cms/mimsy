package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/config"
	"github.com/mimsy-cms/mimsy/internal/storage"
)

type MediaService interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader, contentType string, user *auth.User) (*Media, error)
	GetById(ctx context.Context, id int64) (*Media, error)
	FindAll(ctx context.Context) ([]Media, error)
	GetTemporaryURL(ctx context.Context, media *Media) (string, error)
	Delete(ctx context.Context, media *Media) error
}

type mediaService struct {
	storage         storage.Storage
	mediaRepository Repository
}

func NewService(storage storage.Storage, mediaRepository Repository) MediaService {
	return &mediaService{storage: storage, mediaRepository: mediaRepository}
}

func getBaseAndExt(name string) (string, string) {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	return base, ext
}

func (s *mediaService) resolveFilenameConflict(ctx context.Context, original string) (string, error) {
	name := original
	base, ext := getBaseAndExt(name)
	counter := 1

	for {
		existing, err := s.mediaRepository.FindByName(ctx, name)
		if err != nil {
			return "", err
		}
		if existing == nil {
			break
		}
		name = fmt.Sprintf("%s(%d)%s", base, counter, ext)
		counter++
	}

	return name, nil
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

	var media *Media
	if err := config.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.storage.Upload(ctx, id.String(), file, contentType); err != nil {
			return err
		}

		originalName := fileHeader.Filename
		finalName, err := s.resolveFilenameConflict(ctx, originalName)
		if err != nil {
			return err
		}

		params := &CreateMediaParams{
			Uuid:         id,
			Name:         finalName,
			ContentType:  contentType,
			Size:         fileHeader.Size,
			UploadedById: user.ID,
		}

		media, err = s.mediaRepository.Create(ctx, params)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return media, nil
}

func (s *mediaService) GetById(ctx context.Context, id int64) (*Media, error) {
	return s.mediaRepository.GetById(ctx, id)
}

func (s *mediaService) FindAll(ctx context.Context) ([]Media, error) {
	return s.mediaRepository.FindAll(ctx)
}

func (s *mediaService) GetTemporaryURL(ctx context.Context, media *Media) (string, error) {
	expires := time.Now().Add(3 * time.Hour)

	return s.storage.GetTemporaryURL(media.Uuid.String(), expires)
}

func (s *mediaService) Delete(ctx context.Context, media *Media) error {
	if err := config.WithinTx(ctx, func(ctx context.Context) error {
		if err := s.mediaRepository.Delete(ctx, media); err != nil {
			return err
		}

		if err := s.storage.Delete(ctx, media.Uuid.String()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
