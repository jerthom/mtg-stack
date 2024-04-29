package app

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"mtgStack/identify"
)

func FetchCards(dir string) (map[string]int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir entries: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(entries))

	cardChan := make(chan string, len(entries))

	for _, e := range entries {
		go func(img os.DirEntry) {
			defer wg.Done()
			cardName, err := identify.CardTitle(dir + "/" + img.Name())
			if err != nil {
				fmt.Printf("main: unable to identify card title for img %s: %v", img.Name(), err)
			}
			cardChan <- cardName
		}(e)
	}

	wg.Wait()
	close(cardChan)

	cards := make(map[string]int)
	for card := range cardChan {
		cards[card] = cards[card] + 1
	}

	return cards, nil
}

func RecordCards(outPath string, cards map[string]int) error {
	output, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("main: unable to create output file: %v", err)
	}

	for card, count := range cards {
		_, err = output.WriteString(strconv.Itoa(count) + " " + card + "\n")
		if err != nil {
			return fmt.Errorf("main: unable to write card %s to output file: %v", card, err)
		}
	}

	return nil
}
