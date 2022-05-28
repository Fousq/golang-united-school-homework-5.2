package cache

import "time"

type Cache struct {
	pair        map[string]string
	keyDeadline map[string]*time.Time
}

func NewCache() Cache {
	return Cache{pair: make(map[string]string), keyDeadline: make(map[string]*time.Time)}
}

func isExpired(deadline *time.Time) bool {
	return deadline != nil && time.Now().After(*deadline)
}

func (cache *Cache) removeExpired() {
	for key, deadline := range cache.keyDeadline {
		if isExpired(deadline) {
			delete(cache.keyDeadline, key)
			delete(cache.pair, key)
		}
	}
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.removeExpired()
	for pairKey, value := range cache.pair {
		if pairKey == key {
			return value, true
		}
	}
	return "", false
}

func (cache *Cache) Put(key, value string) {
	cache.removeExpired()
	cache.pair[key] = value
	cache.keyDeadline[key] = nil
}

func (cache *Cache) Keys() []string {
	cache.removeExpired()
	keys := make([]string, 0, len(cache.keyDeadline))
	for k := range cache.keyDeadline {
		keys = append(keys, k)
	}
	return keys
}

func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.removeExpired()
	cache.pair[key] = value
	cache.keyDeadline[key] = &deadline
}
