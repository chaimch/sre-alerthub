package main

import "log"

func main() {
	t := map[int64][]int64{}
	if _, ok := t[1]; !ok {
		log.Println(t[1])
	}
}
