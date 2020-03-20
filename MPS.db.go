package main

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

func storetest() {
	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	d := diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	// Write three bytes to the key "alpha".
	key := "alpha"
	d.Write(key, []byte("test 💩 data"))

	// Read the value back out of the store.
	value, _ := d.Read(key)
	fmt.Printf("%v\n", string(value))

	// Erase the key+value from the store (and the disk).
	d.Erase(key)
}
