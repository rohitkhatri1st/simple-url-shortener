package app

import (
	"math/rand"
	"net/url"
	"simple-url-shortener/model"
	"simple-url-shortener/server/storage"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Url interface {
	ShortenUrl(shortenUrlRequest *model.ShortenUrlRequest) (*model.ShortenUrlResp, error)
	GetOriginalUrlFromShortKey(shortKey string) (string, error)
	GetTopDomains() []model.DomainCount
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

func (ui UrlImpl) GetTopDomains() []model.DomainCount {
	allDomainCounts := ui.Db.DomainCounter.FindAll()

	allDomainCountsSlice := []model.DomainCount{}
	for domain, count := range allDomainCounts {
		allDomainCountsSlice = append(allDomainCountsSlice, model.DomainCount{
			Domain: domain,
			Count:  count.(int),
		})
	}
	sort.Slice(allDomainCountsSlice, func(i, j int) bool {
		return allDomainCountsSlice[i].Count > allDomainCountsSlice[j].Count
	})
	if len(allDomainCountsSlice) < 4 {
		return allDomainCountsSlice
	}
	topDomains := allDomainCountsSlice[:3]
	return topDomains
}

func (ui UrlImpl) GetOriginalUrlFromShortKey(shortKey string) (string, error) {
	urlData, err := ui.Db.UrlData.FindUrlByShortKeyIndex(shortKey)
	if err != nil {
		return "", err
	}
	return urlData, nil
}

func (ui UrlImpl) ShortenUrl(shortenUrlRequest *model.ShortenUrlRequest) (*model.ShortenUrlResp, error) {
	// Adding Mutex so that same url does not get converted twice. It can be removed
	// as it may slow down performance if multiple shortUrls or errors are tolerable.
	ui.Mutex.Lock()
	defer ui.Mutex.Unlock()

	shortKey := ui.generateShortKey()
	shortenUrl := GetServerAddress() + "/" + shortKey

	// Trimming these so that we don't treat "google.com" and "google.com/" as different urls
	originalUrl := strings.Trim(shortenUrlRequest.OriginalUrl, "?")
	originalUrl = strings.Trim(originalUrl, "/")

	originalParsedUrl, err := url.ParseRequestURI(originalUrl)
	if err != nil {
		return nil, err
	}
	domainName := originalParsedUrl.Hostname()
	urlData := model.UrlData{
		ShortenUrl:  shortenUrl,
		OriginalUrl: originalUrl,
		ShortKey:    shortKey,
	}

	urlDataFromDb := &model.UrlData{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		urlDataFromDb = ui.Db.UrlData.InsertUniqueAndGet(originalUrl, urlData)
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
	const keyLength = model.ShortKeyLength
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[ui.RandomGenerator.Intn(len(charset))]
	}
	return string(shortKey)
}
