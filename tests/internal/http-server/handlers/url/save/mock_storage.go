package save

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

func (m *MockURLSaver) SaveURL(ctx context.Context, longUrl string, urlAlias string) error {
	_, exists := m.SavedURLs[urlAlias]
	if exists {
		return storage.ErrUrlAlreadyExists
	}

	m.SavedURLs[urlAlias] = longUrl
	return nil
}

func (m *MockURLSaver) GetURL(ctx context.Context, urlAlias string) (string, error) {
	url, exists := m.SavedURLs[urlAlias]
	if !exists {
		return "", storage.ErrUrlNotFound
	}
	return url, nil
}

func (m *MockURLSaver) DeleteURL(ctx context.Context, urlAlias string) error {
	delete(m.SavedURLs, urlAlias)
	return nil
}

func (m *MockURLSaver) Ping(ctx context.Context) error {
	return nil
}
