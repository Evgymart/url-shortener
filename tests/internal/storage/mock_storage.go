package storage

import (
	"context"
	"urlshort/internal/storage"
)

type MockURLSaver struct {
	SavedURLs map[string]string
}

func NewMockURLSaver() *MockURLSaver {
	return &MockURLSaver{
		SavedURLs: make(map[string]string),
	}
}

func (m *MockURLSaver) SaveURL(_ context.Context, urlAlias string, longUrl string) error {
	_, exists := m.SavedURLs[urlAlias]
	if exists {
		return storage.ErrAliasAlreadyTaken
	}

	m.SavedURLs[urlAlias] = longUrl
	return nil
}

func (m *MockURLSaver) GetURL(_ context.Context, urlAlias string) (string, error) {
	url, exists := m.SavedURLs[urlAlias]
	if !exists {
		return "", storage.ErrUrlNotFound
	}
	return url, nil
}

func (m *MockURLSaver) GetUrl(ctx context.Context, urlAlias string) (string, error) {
	return m.GetURL(ctx, urlAlias)
}

func (m *MockURLSaver) DeleteURL(_ context.Context, urlAlias string) error {
	delete(m.SavedURLs, urlAlias)
	return nil
}

func (m *MockURLSaver) Ping(_ context.Context) error {
	return nil
}
