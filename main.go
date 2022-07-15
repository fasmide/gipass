package main

import (
	"fmt"

	"github.com/fasmide/gipass/store"
)

func main() {
	results, err := store.Logins.Query("facebook")
	if err != nil {
		panic(err)
	}

	for i, v := range results {
		password, err := v.CleartextPassword()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d (%s): '%s'\n", i, v.URL, password)
	}
}
