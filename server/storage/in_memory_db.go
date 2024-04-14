package storage

import (
	"errors"
	"simple-url-shortener/model"
)

type InMemoryDatabases struct {
	UrlDatabase *UrlDatabase
}

type InMemoryDb interface {
	FindAll() map[string]interface{}
	FindUrlByShortKeyIndex(shortKey string) (string, error)
	FindByKey(key string) (*model.UrlData, error)
	Insert(key string, value model.UrlData)
	IncreaseCounter(key string)
	InsertUniqueAndGet(key string, value model.UrlData) *model.UrlData
}

func InitDb() *InMemoryDatabases {
	return &InMemoryDatabases{
		UrlDatabase: &UrlDatabase{
			UrlData:       NewInMemoryDb(),
			DomainCounter: NewInMemoryDb(),
		},
	}
}

type UrlDatabase struct {
	UrlData       InMemoryDb
	DomainCounter InMemoryDb
}

type InMemoryDbImpl struct {
	Data                       map[string]interface{}
	ShortKeyToOriginalUrlIndex map[string]string
}

func NewInMemoryDb() InMemoryDb {
	inMemoryDb := InMemoryDbImpl{
		Data:                       map[string]interface{}{},
		ShortKeyToOriginalUrlIndex: map[string]string{},
	}
	return inMemoryDb
}

func (imd InMemoryDbImpl) FindAll() map[string]interface{} {
	return imd.Data
}

func (imd InMemoryDbImpl) FindUrlByShortKeyIndex(shortKey string) (string, error) {
	indexValue, found := imd.ShortKeyToOriginalUrlIndex[shortKey]
	if !found {
		return "", errors.New("index not found")
	}
	// We can return below as well if entire url data is required
	// return imd.FindByKey(indexValue)

	return indexValue, nil
}

// Key is the originalUrl
func (imd InMemoryDbImpl) FindByKey(key string) (*model.UrlData, error) {
	value, found := imd.Data[key]
	if !found {
		return nil, errors.New("key not found")
	}
	urlData := value.(model.UrlData)
	return &urlData, nil
}

func (imd InMemoryDbImpl) Insert(key string, value model.UrlData) {
	imd.Data[key] = value
	imd.ShortKeyToOriginalUrlIndex[value.ShortKey] = key
}

// This method is only for DomainCounterDb
func (imd InMemoryDbImpl) IncreaseCounter(key string) {
	if val, found := imd.Data[key]; found {
		imd.Data[key] = val.(int) + 1
		return
	}
	imd.Data[key] = int(1)
}

func (imd InMemoryDbImpl) InsertUniqueAndGet(key string, value model.UrlData) *model.UrlData {
	if val, found := imd.Data[key]; found {
		urlData := val.(model.UrlData)
		return &urlData
	}
	imd.Data[key] = value
	imd.ShortKeyToOriginalUrlIndex[value.ShortKey] = key
	return &value
}
