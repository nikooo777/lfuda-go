package lfuda

import (
	"sync"

	"github.com/bparli/lfuda-go/simplelfuda"
)

// Cache is a thread-safe fixed size lfuda cache.
type Cache struct {
	lfuda simplelfuda.LFUDACache
	lock  sync.RWMutex
}

// New creates an lfuda of the given size.
func New(size float64) *Cache {
	return newWithEvict(size, "LFUDA", nil)
}

// NewGDSF creates an lfuda of the given size and the GDSF cache policy.
func NewGDSF(size float64) *Cache {
	return newWithEvict(size, "GDSF", nil)
}

// NewLFU creates an lfuda of the given size.
func NewLFU(size float64) *Cache {
	return newWithEvict(size, "LFU", nil)
}

// NewWithEvict constructs a fixed size LFUDA cache with the given eviction
// callback.
func NewWithEvict(size float64, onEvicted func(key interface{}, value interface{})) *Cache {
	return newWithEvict(size, "LFUDA", onEvicted)
}

// NewGDSFWithEvict constructs a fixed GDSF size cache with the given eviction
// callback.
func NewGDSFWithEvict(size float64, onEvicted func(key interface{}, value interface{})) *Cache {
	return newWithEvict(size, "GDSF", onEvicted)
}

// NewLFUWithEvict constructs a fixed size LFU cache with the given eviction
// callback.
func NewLFUWithEvict(size float64, onEvicted func(key interface{}, value interface{})) *Cache {
	return newWithEvict(size, "LFU", onEvicted)
}

func newWithEvict(size float64, policy string, onEvicted func(key interface{}, value interface{})) *Cache {
	if policy == "GDSF" {
		gdsf := simplelfuda.NewGDSF(size, simplelfuda.EvictCallback(onEvicted))
		return &Cache{
			lfuda: gdsf,
		}
	} else if policy == "LFU" {
		lfu := simplelfuda.NewLFU(size, simplelfuda.EvictCallback(onEvicted))
		return &Cache{
			lfuda: lfu,
		}
	}
	lfuda := simplelfuda.NewLFUDA(size, simplelfuda.EvictCallback(onEvicted))
	return &Cache{
		lfuda: lfuda,
	}
}

// Purge is used to completely clear the cache.
func (c *Cache) Purge() {
	c.lock.Lock()
	c.lfuda.Purge()
	c.lock.Unlock()
}

// Set adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) Set(key, value interface{}) (ok bool) {
	c.lock.Lock()
	ok = c.lfuda.Set(key, value)
	c.lock.Unlock()
	return ok
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lfuda.Get(key)
	c.lock.Unlock()
	return value, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *Cache) Contains(key interface{}) bool {
	c.lock.RLock()
	containKey := c.lfuda.Contains(key)
	c.lock.RUnlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	value, ok = c.lfuda.Peek(key)
	c.lock.RUnlock()
	return value, ok
}

// ContainsOrSet checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether the key/value was set or not.
func (c *Cache) ContainsOrSet(key, value interface{}) (ok, set bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lfuda.Contains(key) {
		return true, false
	}
	set = c.lfuda.Set(key, value)
	return false, set
}

// PeekOrSet checks if a key is in the cache without updating the
// hits or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether the key/value was set or not.
func (c *Cache) PeekOrSet(key, value interface{}) (previous interface{}, ok, set bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	previous, ok = c.lfuda.Peek(key)
	if ok {
		return previous, true, false
	}

	set = c.lfuda.Set(key, value)
	return nil, false, set
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key interface{}) (present bool) {
	c.lock.Lock()
	present = c.lfuda.Remove(key)
	c.lock.Unlock()
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lfuda.Keys()
	c.lock.RUnlock()
	return keys
}

// Len returns the number of items in the cache.
func (c *Cache) Len() (length int) {
	c.lock.RLock()
	length = c.lfuda.Len()
	c.lock.RUnlock()
	return length
}

// Size returns the current size of the cache in bytes.
func (c *Cache) Size() (size float64) {
	c.lock.RLock()
	size = c.lfuda.Size()
	c.lock.RUnlock()
	return size
}

// Age returns the cache's current age
func (c *Cache) Age() (age float64) {
	c.lock.RLock()
	age = c.lfuda.Age()
	c.lock.RUnlock()
	return age
}
