package storage

import (
	"context"
	"io"
	"log/slog"

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
}

type swiftStorage struct {
	connection swift.Connection
	container  string
}

type SwiftOptionFn func(*swiftOption)

func NewSwift(opts ...SwiftOptionFn) *swiftStorage {
	o := &swiftOption{}

	for _, opt := range opts {
		opt(o)
	}

	return &swiftStorage{
		container: o.container,
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

func (s *swiftStorage) Upload(ctx context.Context, id string, data io.Reader, contentType string) error {
	if _, err := s.connection.ObjectPut(ctx, s.container, id, data, false, "", contentType, nil); err != nil {
		slog.Error("Failed to upload file to swift", "error", err)
		return err
	}

	return nil
}
