package storage

import (
	"context"
	"io"
)

type Storage interface {
	// Upload uploads a file to the storage backend.
	//
	// It returns an error if the upload fails or if a file with the same id already exists.
	Upload(ctx context.Context, id string, data io.Reader, contentType string) error
}
