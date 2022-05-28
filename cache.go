package cache

import "time"

type CachePair struct {
	key       string
	value     string
	expiredIn *time.Time
}

type Cache struct {
	pair []CachePair
}

func NewCache() Cache {
	return Cache{}
}

func (cachePair *CachePair) isExpired() bool {
	return cachePair.expiredIn != nil && time.Now().After(*cachePair.expiredIn)
}

func (cache *Cache) removeExpired() {
	expiredIndexes := make([]int, 0)
	for i := 0; i < len(cache.pair); i++ {
		if cache.pair[i].isExpired() {
			expiredIndexes = append(expiredIndexes, i)
		}
	}
	for _, expiredPairIndex := range expiredIndexes {
		cache.pair = append(cache.pair[:expiredPairIndex], cache.pair[expiredPairIndex+1:]...)
	}
}

func (cache *Cache) Get(key string) (string, bool) {
	cache.removeExpired()
	for _, cachePair := range cache.pair {
		if cachePair.key == key && !cachePair.isExpired() {
			return cachePair.value, true
		}
	}
	return "", false
}

func (cache *Cache) Put(key, value string) {
	cache.removeExpired()
	cache.pair = append(cache.pair, CachePair{key, value, nil})
}

func (cache *Cache) Keys() []string {
	cache.removeExpired()
	keys := make([]string, 0, len(cache.pair))
	for _, cachePair := range cache.pair {
		if !cachePair.isExpired() {
			keys = append(keys, cachePair.key)
		}
	}
	return keys
}

func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.removeExpired()
	cache.pair = append(cache.pair, CachePair{key, value, &deadline})
}
