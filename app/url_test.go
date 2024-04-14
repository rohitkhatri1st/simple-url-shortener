package app

import (
	"simple-url-shortener/model"
	"simple-url-shortener/server/mocks"
	"simple-url-shortener/server/storage"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
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

// This is not full functional test case but partial without setting up server.
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

func TestUrlImpl_GetOriginalUrlFromShortKeyUnitTest(t *testing.T) {
	mockedDatabase := &mocks.InMemoryDb{}
	mockedDatabase.On("FindUrlByShortKeyIndex", mock.AnythingOfType("string")).
		Return("http://example.com", nil).
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

	expectedOriginalUrl := "http://example.com"
	result, err := ui.GetOriginalUrlFromShortKey("iYIUZj")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedOriginalUrl {
		t.Errorf("Expected: %s, Got: %s", expectedOriginalUrl, result)
	}
}

// This is not full functional test case but partial without setting up server.
func TestUrlImpl_GetOriginalUrlFromShortKeyFunctionalTest(t *testing.T) {
	db := storage.InitDb().UrlDatabase
	ui := InitUrl(&UrlImplOpts{
		App: &App{
			Logger: &zerolog.Logger{},
		},
		Db: db,
	})

	expectedOriginalUrl := "http://example.com"
	shortenUrlResp, err := ui.ShortenUrl(&model.ShortenUrlRequest{OriginalUrl: expectedOriginalUrl})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	shortKey := shortenUrlResp.ShortenUrl[len(GetServerAddress())+1:]
	result, err := ui.GetOriginalUrlFromShortKey(shortKey)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedOriginalUrl {
		t.Errorf("Expected: %s, Got: %s", expectedOriginalUrl, result)
	}
}

// This is not full functional test case but partial without setting up server.
func TestUrlImpl_GetTopDomainsFunctionalTest(t *testing.T) {
	db := storage.InitDb().UrlDatabase
	ui := InitUrl(&UrlImplOpts{
		App: &App{
			Logger: &zerolog.Logger{},
		},
		Db: db,
	})

	urls := []string{
		"http://example.com/path1",
		"http://example.com/path2",
		"http://google.com",
		"http://yahoo.com",
		"http://yahoo.com/path1",
		"http://bing.com",
		"http://bing.com/path1",
		"http://bing.com/path2",
	}
	for _, url := range urls {
		_, err := ui.ShortenUrl(&model.ShortenUrlRequest{OriginalUrl: url})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	topDomainsResult := ui.GetTopDomains()

	expectedTopDomains := []model.DomainCount{
		{Domain: "bing.com", Count: 3},
		{Domain: "example.com", Count: 2},
		{Domain: "yahoo.com", Count: 2},
	}
	assert.Equal(t, len(expectedTopDomains), len(topDomainsResult))

	for i, expected := range expectedTopDomains {
		assert.Equal(t, expected.Domain, topDomainsResult[i].Domain)
		assert.Equal(t, expected.Count, topDomainsResult[i].Count)
	}

	// assert.True(t, reflect.DeepEqual(expectedTopDomains, topDomainsResult))
}
