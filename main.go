package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"
)

type translation map[string]string

func measure() {
	var t translation

	r, err := http.Get("https://raw.githubusercontent.com/awdarlan/polytopia/master/kk_KZ.json")
	if err != nil {
		panic(err)
	}

	if err = json.NewDecoder(r.Body).Decode(&t); err!=nil {
		panic(err)
	}

	var tred, cnt uint64

	for _,v := range t {
		for _, r := range v {
			if unicode.Is(unicode.Cyrillic, r) {
				tred++
				break
			}
		}
		cnt++
	}

	fmt.Printf("\n\n\nсөздіктің %d%% аударылған, яғни %d/%d\n\n\n", int(float64(tred)/float64(cnt) * 100), tred, cnt)
}

func main() {
	measure()
}
