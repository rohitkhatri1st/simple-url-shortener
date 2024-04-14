package storage

type InMemoryDb map[string]interface{}

func NewInMemoryDb() map[string]interface{} {
	db := map[string]interface{}{}
	return InMemoryDb(db)
}
