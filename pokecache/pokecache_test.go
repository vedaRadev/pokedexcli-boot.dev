package pokecache

import(
    "testing"
    "time"
    "fmt"
)

func TestAddGet(t *testing.T) {
    const interval = 0 // no reaping
    cases := []struct { key string; val []byte } {
        { key: "http://example.com", val: []byte("test data") },
        { key: "http://example.com/path", val: []byte("test data at path") },
        { key: "random test", val: []byte("abc123") },
    }

    for i, c := range cases {
        t.Run(fmt.Sprintf("Test %v", i), func (t *testing.T) {
            cache := NewCache(interval)
            cache.Add(c.key, c.val)
            val, ok := cache.Get(c.key)
            if !ok {
                t.Errorf("key not found")
                return
            }

            if string(val) != string(c.val) {
                t.Errorf("retrieved an unexpected value from the cache")
                return
            }
        })
    }
}

func TestReaping(t *testing.T) {
    const interval = 5 * time.Millisecond
    const wait = interval + 5 * time.Millisecond
    cache := NewCache(interval)
    cache.Add("test", []byte("test data"))

    _, ok := cache.Get("test")
    if !ok {
        t.Errorf("entry was either not inserted or was reaped too soon!")
        return
    }

    time.Sleep(wait)
    _, ok = cache.Get("test")
    if ok {
        t.Errorf("expected entry to be reaped but it was still there")
        return
    }
}
