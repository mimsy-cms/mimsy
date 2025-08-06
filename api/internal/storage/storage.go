package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	// Authenticate authenticates the storage backend.
	// It returns an error if the authentication fails.
	Authenticate(ctx context.Context) error
	// Upload uploads a file to the storage backend.
	//
	// It returns an error if the upload fails or if a file with the same id already exists.
	Upload(ctx context.Context, id string, data io.Reader, contentType string) error
	// Delete deletes a file from the storage backend.
	Delete(ctx context.Context, id string) error
	// GetTemporaryURL returns a temporary URL for the file with the given id.
	GetTemporaryURL(id string, expires time.Time) (string, error)
}
