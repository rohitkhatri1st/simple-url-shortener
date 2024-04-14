package storage

import "errors"

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

type InMemoryDb map[string]interface{}

func NewInMemoryDb() map[string]interface{} {
	db := map[string]interface{}{}
	return InMemoryDb(db)
}

func (imd InMemoryDb) Find(key string) (interface{}, error) {
	value, found := imd[key]
	if !found {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (imd InMemoryDb) Insert(key string, value interface{}) {
	imd[key] = value
}

func (imd InMemoryDb) IncreaseCounter(key string) {
	if val, found := imd[key]; found {
		imd[key] = val.(int) + 1
		return
	}
	imd[key] = int(1)
}

func (imd InMemoryDb) InsertUniqueAndGet(key string, value interface{}) interface{} {
	if val, found := imd[key]; found {
		return val
	}
	imd[key] = value
	return value
}
