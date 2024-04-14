package app

import (
	"simple-url-shortener/model"
	"simple-url-shortener/server/mocks"
	"simple-url-shortener/server/storage"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestUrlImpl_ShortenUrlUnitTest(t *testing.T) {
	originalUrl := "http://example.com"
	mockedDatabase := &mocks.InMemoryDb{}
	mockedDatabase.On("InsertUniqueAndGet", mock.AnythingOfType("string"), mock.AnythingOfType("model.UrlData")).
		Return(&model.UrlData{
			OriginalUrl: originalUrl,
			ShortenUrl:  GetServerAddress() + "/iYIUZj",
			ShortKey:    "iYIUZj",
		}).
		Once()
	mockedDatabase.On("IncreaseCounter", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	ui := InitUrl(&UrlImplOpts{
		App: &App{
			Logger: &zerolog.Logger{},
		},
		Db: &storage.UrlDatabase{
			UrlData:       mockedDatabase,
			DomainCounter: mockedDatabase,
		},
	})

	shortenUrlRequest := &model.ShortenUrlRequest{OriginalUrl: originalUrl}

	result, err := ui.ShortenUrl(shortenUrlRequest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedOriginalUrl := "http://example.com"
	expectedUrlData := model.UrlData{
		OriginalUrl: originalUrl,
		ShortenUrl:  GetServerAddress() + "/iYIUZj",
		ShortKey:    "iYIUZj",
	}
	expectedShortenUrl := expectedUrlData.ShortenUrl
	if result.OriginalUrl != expectedOriginalUrl {
		t.Errorf("Expected OriginalUrl: %s, Got: %s", expectedOriginalUrl, result.OriginalUrl)
	}
	if result.ShortenUrl != expectedShortenUrl {
		t.Errorf("Expected ShortenUrl: %s, Got: %s", expectedShortenUrl, result.ShortenUrl)
	}
	if !strings.HasPrefix(result.ShortenUrl, expectedShortenUrl) {
		t.Errorf("Expected ShortenUrl to start with: %s, Got: %s", expectedShortenUrl, result.ShortenUrl)
	}
	if len(result.ShortenUrl[len(GetServerAddress()+"/"):]) != model.ShortKeyLength {
		t.Errorf("Expected ShortenUrl length is: %d, Got: %d", model.ShortKeyLength, len(result.ShortenUrl[len(GetServerAddress()+"/"):]))
	}
}

func TestUrlImpl_ShortenUrlFunctionalTest(t *testing.T) {
	db := storage.InitDb().UrlDatabase
	ui := InitUrl(&UrlImplOpts{
		App: &App{
			Logger: &zerolog.Logger{},
		},
		Db: db,
	})

	shortenUrlRequest := &model.ShortenUrlRequest{OriginalUrl: "http://example.com"}

	result, err := ui.ShortenUrl(shortenUrlRequest)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedOriginalUrl := "http://example.com"
	urlData, err := db.UrlData.FindByKey(expectedOriginalUrl)
	if err != nil {
		t.Error(err)
	}
	expectedShortenUrl := urlData.ShortenUrl
	if result.OriginalUrl != expectedOriginalUrl {
		t.Errorf("Expected OriginalUrl: %s, Got: %s", expectedOriginalUrl, result.OriginalUrl)
	}
	if result.ShortenUrl != expectedShortenUrl {
		t.Errorf("Expected ShortenUrl: %s, Got: %s", expectedShortenUrl, result.ShortenUrl)
	}
	if !strings.HasPrefix(result.ShortenUrl, expectedShortenUrl) {
		t.Errorf("Expected ShortenUrl to start with: %s, Got: %s", expectedShortenUrl, result.ShortenUrl)
	}
	if len(result.ShortenUrl[len(GetServerAddress()+"/"):]) != model.ShortKeyLength {
		t.Errorf("Expected ShortenUrl length is: %d, Got: %d", model.ShortKeyLength, len(result.ShortenUrl[len(GetServerAddress()+"/"):]))
	}
}

func TestUrlImpl_ShortenUrl_InvalidUrl(t *testing.T) {
	db := storage.InitDb().UrlDatabase
	ui := InitUrl(&UrlImplOpts{
		App: &App{
			Logger: &zerolog.Logger{},
		},
		Db: db,
	})

	shortenUrlRequest := &model.ShortenUrlRequest{OriginalUrl: "123"}

	result, err := ui.ShortenUrl(shortenUrlRequest)
	if err == nil || result != nil {
		t.Errorf("Error was expected but not found")
	}
}
