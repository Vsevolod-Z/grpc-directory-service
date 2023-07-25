package cache

import (
	"container/heap"
	pb "grpc-directory-service/api/directory"
	"log"
	"sync"
	"time"
)

type CachedData struct {
	Dir        string
	Files      []*pb.FileInfo
	InsertTime time.Time
	HitCount   int
}
type Cache struct {
	mu         sync.RWMutex
	cache      map[string]*CachedData
	expiration time.Duration
	capacity   int
	leastUsed  cacheHeap
}
type cacheHeap []*CachedData

func (h cacheHeap) Len() int           { return len(h) }
func (h cacheHeap) Less(i, j int) bool { return h[i].HitCount < h[j].HitCount }
func (h cacheHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *cacheHeap) Push(x interface{}) {
	*h = append(*h, x.(*CachedData))
}

func (h *cacheHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewCache(expiration time.Duration, capacity int) *Cache {
	c := &Cache{
		cache:      make(map[string]*CachedData),
		expiration: expiration,
		capacity:   capacity,
	}
	heap.Init(&c.leastUsed)
	log.Printf("the cache has been created\n expiration: %s,capacity:   %d", c.expiration, c.capacity)
	return c
}
func (c *Cache) Get(key string) ([]*pb.FileInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.cache[key]
	if !found {
		log.Printf("%s | was not found in the cache", key)
		go c.cleanup()
		return nil, false
	}

	entry.HitCount++
	if len(c.leastUsed) > 2 {
		heap.Fix(&c.leastUsed, entry.HitCount-1)
	}
	log.Printf("%s | is taken from the cache", key)
	go c.cleanup()
	return entry.Files, true
}
func (c *Cache) Add(key string, files []*pb.FileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.cache) >= c.capacity {
		log.Printf("cache is full , deleting least Used Entry: %s", key)
		leastUsedEntry := heap.Pop(&c.leastUsed).(*CachedData)
		delete(c.cache, leastUsedEntry.Dir)
	}

	entry := &CachedData{
		Dir:        key,
		Files:      files,
		InsertTime: time.Now(),
		HitCount:   1,
	}
	c.cache[key] = entry
	heap.Push(&c.leastUsed, entry)
	log.Printf("key: %s | was added to the cache", key)

}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	counter := 0
	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.InsertTime) >= c.expiration {
			delete(c.cache, key)
			counter++
		}
	}
	log.Printf("cache cleaned, %d entries deleted", counter)
}
