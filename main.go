package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"mtgStack/identify"
)

func main() {
	entries, err := os.ReadDir("./img")
	if err != nil {
		log.Fatalf("Unable to read dir entries: %v", err)
	}

	cards := make(map[string]int)

	for _, e := range entries {
		cardName, err := identify.CardTitle("img/" + e.Name())
		if err != nil {
			fmt.Printf("main: unable to identify card title for img %s: %v", e.Name(), err)
		}
		cards[cardName] = cards[cardName] + 1
	}

	output, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("main: unable to create output file: %v", err)
	}

	for card, count := range cards {
		_, err = output.WriteString(strconv.Itoa(count) + " " + card + "\n")
		if err != nil {
			log.Printf("main: unable to write card %s to output file: %v", card, err)
		}
	}
}
