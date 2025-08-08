package storage

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/ncw/swift/v2"
)

type swiftOption struct {
	username  string
	apiKey    string
	authURL   string
	domain    string
	tenant    string
	container string
	region    string
	secretKey string
}

type swiftStorage struct {
	connection swift.Connection
	container  string
	// secretKey is used to sign the temporary URLs.
	secretKey string
}

type SwiftOptionFn func(*swiftOption)

func NewSwift(opts ...SwiftOptionFn) *swiftStorage {
	o := &swiftOption{}

	for _, opt := range opts {
		opt(o)
	}

	return &swiftStorage{
		container: o.container,
		secretKey: o.secretKey,
		connection: swift.Connection{
			UserName: o.username,
			ApiKey:   o.apiKey,
			AuthUrl:  o.authURL,
			Domain:   o.domain,
			Tenant:   o.tenant,
			Region:   o.region,
		},
	}
}

func WithSwiftUsername(username string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.username = username
	}
}

func WithSwiftApiKey(apiKey string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.apiKey = apiKey
	}
}

func WithSwiftAuthURL(authURL string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.authURL = authURL
	}
}

func WithSwiftDomain(domain string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.domain = domain
	}
}

func WithSwiftTenant(tenant string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.tenant = tenant
	}
}

func WithSwiftContainer(container string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.container = container
	}
}

func WithSwiftRegion(region string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.region = region
	}
}

func WithSwiftSecretKey(secretKey string) SwiftOptionFn {
	return func(s *swiftOption) {
		s.secretKey = secretKey
	}
}

func (s *swiftStorage) Authenticate(ctx context.Context) error {
	return s.connection.Authenticate(ctx)
}

func (s *swiftStorage) Upload(ctx context.Context, id string, data io.Reader, contentType string) error {
	if _, err := s.connection.ObjectPut(ctx, s.container, id, data, false, "", contentType, nil); err != nil {
		slog.Error("Failed to upload file to swift", "error", err)
		return err
	}

	return nil
}

func (s *swiftStorage) GetTemporaryURL(id string, expires time.Time) (string, error) {
	tempURL := s.connection.ObjectTempUrl(s.container, id, s.secretKey, "GET", expires)

	return tempURL, nil
}

func (s *swiftStorage) Delete(ctx context.Context, id string) error {
	if err := s.connection.ObjectDelete(ctx, s.container, id); err != nil {
		slog.Error("Failed to delete file from swift", "error", err)
		return err
	}

	return nil
}
