package storage

import (
	"errors"
	"simple-url-shortener/model"
)

type InMemoryDatabases struct {
	UrlDatabase *UrlDatabase
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

type InMemoryDb struct {
	Data                       map[string]interface{}
	ShortKeyToOriginalUrlIndex map[string]string
}

func NewInMemoryDb() InMemoryDb {
	inMemoryDb := InMemoryDb{
		Data:                       map[string]interface{}{},
		ShortKeyToOriginalUrlIndex: map[string]string{},
	}
	return inMemoryDb
}

func (imd InMemoryDb) FindUrlByShortKeyIndex(shortKey string) (string, error) {
	indexValue, found := imd.ShortKeyToOriginalUrlIndex[shortKey]
	if !found {
		return "", errors.New("index not found")
	}
	// We can return below as well if entire url data is required
	// return imd.FindByKey(indexValue)

	return indexValue, nil
}

func (imd InMemoryDb) FindByKey(key string) (*model.UrlData, error) {
	value, found := imd.Data[key]
	if !found {
		return nil, errors.New("key not found")
	}
	urlData := value.(model.UrlData)
	return &urlData, nil
}

func (imd InMemoryDb) Insert(key string, value model.UrlData) {
	imd.Data[key] = value
	imd.ShortKeyToOriginalUrlIndex[value.ShortKey] = key
}

func (imd InMemoryDb) IncreaseCounter(key string) {
	if val, found := imd.Data[key]; found {
		imd.Data[key] = val.(int) + 1
		return
	}
	imd.Data[key] = int(1)
}

func (imd InMemoryDb) InsertUniqueAndGet(key string, value model.UrlData) *model.UrlData {
	if val, found := imd.Data[key]; found {
		urlData := val.(model.UrlData)
		return &urlData
	}
	imd.Data[key] = value
	imd.ShortKeyToOriginalUrlIndex[value.ShortKey] = key
	return &value
}
