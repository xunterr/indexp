package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/xunterr/indexp/cmd"
)

func main() {
	flag.Parse()
	cmd.Execute()

}

func DisplayFT(ft map[string]float64, num int) {
	keys := make([]string, 0, len(ft))

	for key := range ft {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return ft[keys[i]] > ft[keys[j]]
	})

	for i := 0; i < num && i < len(keys); i++ {
		key := keys[i]
		fmt.Printf("%d. %s -- %f \n", i+1, key, ft[key])
	}
}
