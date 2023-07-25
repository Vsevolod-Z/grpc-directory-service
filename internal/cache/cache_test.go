package cache

import (
	pb "grpc-directory-service/api/directory"
	"testing"
	"time"
)

type TestCache struct {
	*Cache
}

func NewTestCache(expiration time.Duration, capacity int) *TestCache {
	return &TestCache{
		Cache: NewCache(expiration, capacity),
	}
}

func TestNewCache(t *testing.T) {
	expiration := 30 & time.Second
	capacity := 100

	c := NewTestCache(expiration, capacity)
	//Cache not nil
	if c == nil {
		t.Errorf("Expected a non-nil cache, got nil")
	}
	//Empty heap when creating
	if len(c.leastUsed) != 0 {
		t.Errorf("Expected cache to be empty, but it has %d entries", len(c.leastUsed))
	}
	//Matching expiration
	if c.expiration != expiration {
		t.Errorf("Expected expiration time to be %v, but got %v", expiration, c.expiration)
	}
	//Matching capacity
	if c.capacity != capacity {
		t.Errorf("Expected capacity to be %d, but got %d", capacity, c.capacity)
	}
}

func TestAddAndGet(t *testing.T) {
	expiration := 10 & time.Second
	capacity := 100
	fileInfo1 := &pb.FileInfo{Name: "QwёЯ@日本Пр123안녕34", Size: 9213372036854775807}
	fileInfo2 := &pb.FileInfo{Name: "Ü×ч67فغHello반사!", Size: 9203372036854775807}

	c := NewTestCache(expiration, capacity)
	c.Add("key1", []*pb.FileInfo{fileInfo1, fileInfo2})
	c.Add("key2", []*pb.FileInfo{fileInfo2, fileInfo1})

	files, found := c.Get("key1")
	//existing key
	if !found {
		t.Errorf("Expected found = %v, but got %v", true, found)
	}
	//Matching name
	if fileInfo1.Name != files[0].Name {
		t.Errorf("Expected Name = %v, but got %v", fileInfo1.Name, files[0].Name)
	}
	//Matching size
	if fileInfo1.Size != files[0].Size {
		t.Errorf("Expected Size = %v, but got %v", fileInfo1.Size, files[0].Size)
	}
	// Wait for cache entry to expire
	time.Sleep(10 * time.Second)
	// non-existing key
	files, found = c.Get("key3")
	if found {
		t.Errorf("Expected found = %v, but got %v", false, found)
	}
	if files != nil {
		t.Errorf("Expected files = %v, but got %v", nil, files)
	}
	// Test getting expired key
	files, found = c.Get("key2")
	if found {
		t.Errorf("Expected found = %v, but got %v", false, found)
	}
	if files != nil {
		t.Errorf("Expected files = %v, but got %v", nil, files)
	}

}
