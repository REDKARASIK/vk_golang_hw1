package main

import (
	"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	mkey := "habrTag"

	memcacheClient.Set(&memcache.Item{
		Key:        mkey,
		Value:      []byte("1"),
		Expiration: 3,
	})

	_, _ = memcacheClient.Increment(mkey, 1)

	item, err := memcacheClient.Get(mkey)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		fmt.Println("MC error", err)
	}

	fmt.Printf("mc value %#v\n", string(item.Value))

	_ = memcacheClient.Delete(mkey)

	item, err = memcacheClient.Get(mkey)
	if errors.Is(err, memcache.ErrCacheMiss) {
		fmt.Println("record not found in MC")
	}
}
