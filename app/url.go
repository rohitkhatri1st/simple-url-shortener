package app

import (
	"math/rand"
	"net/url"
	"simple-url-shortener/model"
	"simple-url-shortener/server/storage"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Url interface {
	ShortenUrl(shortenUrlRequest *model.ShortenUrlRequest) (*model.ShortenUrlResp, error)
}

type UrlImplOpts struct {
	App    *App
	Db     *storage.UrlDatabase
	Logger *zerolog.Logger
}
type UrlImpl struct {
	App             *App
	Db              *storage.UrlDatabase
	Logger          *zerolog.Logger
	RandomGenerator *rand.Rand
	Mutex           *sync.Mutex
}

// InitUrl returns new instance of url implementation
func InitUrl(opts *UrlImplOpts) Url {
	l := opts.App.Logger.With().Str("service", "url").Logger()
	ui := UrlImpl{
		App:             opts.App,
		Db:              opts.Db,
		Logger:          &l,
		Mutex:           &sync.Mutex{},
		RandomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &ui
}

func (ui UrlImpl) ShortenUrl(shortenUrlRequest *model.ShortenUrlRequest) (*model.ShortenUrlResp, error) {
	ui.Mutex.Lock()
	defer ui.Mutex.Unlock()

	shortKey := ui.generateShortKey()
	shortenUrl := "http://localhost:8000/" + shortKey

	// Trimming these so that we don't treat "google.com" and "google.com/" as different urls
	originalUrl := strings.Trim(shortenUrlRequest.OriginalUrl, "?")
	originalUrl = strings.Trim(originalUrl, "/")

	originalParsedUrl, err := url.Parse(originalUrl)
	if err != nil {
		return nil, err
	}
	domainName := originalParsedUrl.Hostname()
	urlData := model.UrlData{
		ShortenUrl:  shortenUrl,
		OriginalUrl: originalUrl,
	}

	urlDataFromDb := model.UrlData{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		urlDataFromDb = ui.Db.UrlData.InsertUniqueAndGet(shortenUrlRequest.OriginalUrl, urlData).(model.UrlData)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ui.Db.DomainCounter.IncreaseCounter(domainName)
	}()
	wg.Wait()
	shortenUrlResp := &model.ShortenUrlResp{
		OriginalUrl: urlDataFromDb.OriginalUrl,
		ShortenUrl:  urlDataFromDb.ShortenUrl,
	}
	return shortenUrlResp, nil
}

func (ui UrlImpl) generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[ui.RandomGenerator.Intn(len(charset))]
	}
	return string(shortKey)
}
