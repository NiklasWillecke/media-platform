package cache

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewLRUCache(12)                        //Bytes
	cache.Set("1", []byte{'W', 'o', 'r', 'l', 'd'}) //5

	cache.Set("2", []byte{'T', 'e', 's', 't'}) //4

	cache.Set("3", []byte{'H', 'e', 'l', 'l', 'o'}) //5

	value, ok := cache.Get("2")
	fmt.Println(ok)
	if !ok && bytes.Equal(value, []byte{'W', 'o', 'r', 'l', 'd'}) {
		t.Error()
	}
}
