package github_fetcher

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v74/github"
)

type installationToken struct {
	token     string
	expiresAt time.Time
}

type authManager struct {
	appID             int64
	privateKey        *rsa.PrivateKey
	installationCache map[int64]*installationToken
	cacheMutex        sync.RWMutex
}

func newAuthManager(appID int64, privateKeyPEM []byte) (*authManager, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &authManager{
		appID:             appID,
		privateKey:        privateKey,
		installationCache: make(map[int64]*installationToken),
	}, nil
}

func (am *authManager) generateJWT() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": jwt.NewNumericDate(now.Add(-60 * time.Second)),
		"exp": jwt.NewNumericDate(now.Add(5 * time.Minute)),
		"iss": am.appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(am.privateKey)
}

func (am *authManager) createJWTClient(context.Context) (*http.Client, error) {
	jwtToken, err := am.generateJWT()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &http.Client{
		Transport: &jwtTransport{
			token: jwtToken,
			base:  http.DefaultTransport,
		},
	}, nil
}

func (am *authManager) getInstallationID(ctx context.Context, owner, repo string) (int64, error) {
	jwtClient, err := am.createJWTClient(ctx)
	if err != nil {
		return 0, err
	}

	client := github.NewClient(jwtClient)
	installation, _, err := client.Apps.FindRepositoryInstallation(ctx, owner, repo)
	if err != nil {
		return 0, fmt.Errorf("failed to find installation for %s/%s: %w", owner, repo, err)
	}

	return installation.GetID(), nil
}

func (am *authManager) getInstallationIDWithClient(ctx context.Context, client *github.Client, owner, repo string) (int64, error) {
	installation, _, err := client.Apps.FindRepositoryInstallation(ctx, owner, repo)
	if err != nil {
		return 0, fmt.Errorf("failed to find installation for %s/%s: %w", owner, repo, err)
	}

	return installation.GetID(), nil
}

func (am *authManager) getInstallationToken(ctx context.Context, installationID int64) (string, error) {
	am.cacheMutex.RLock()
	cached, exists := am.installationCache[installationID]
	am.cacheMutex.RUnlock()

	if exists && time.Now().Before(cached.expiresAt.Add(-5*time.Minute)) {
		return cached.token, nil
	}

	jwtClient, err := am.createJWTClient(ctx)
	if err != nil {
		return "", err
	}

	client := github.NewClient(jwtClient)
	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create installation token: %w", err)
	}

	am.cacheMutex.Lock()
	am.installationCache[installationID] = &installationToken{
		token:     token.GetToken(),
		expiresAt: token.GetExpiresAt().Time,
	}
	am.cacheMutex.Unlock()

	return token.GetToken(), nil
}

type jwtTransport struct {
	token string
	base  http.RoundTripper
}

func (t *jwtTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return t.base.RoundTrip(req)
}
