package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = time.Second * 5
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "www.example.com",
			val: []byte("exampledata"),
		},
		{
			key: "www.example.com/path",
			val: []byte("moreexampledata"),
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("testAddGet test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			cval, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("Expected to find key")
			}
			if string(cval) != string(c.val) {
				t.Errorf("expected to find value")
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
