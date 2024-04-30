package identify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

// CardTitle returns the assumed title for the image provided
// the title is defined as the text in the image that precedes the first newline
func CardTitle(imgPath string) (string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to create client: %v", err)
	}
	defer client.Close()

	file, err := os.Open(imgPath)
	if err != nil {
		return "", fmt.Errorf("unable to open image: %v", err)
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return "", fmt.Errorf("unable to get image from reader: %v", err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return "", fmt.Errorf("unable to get text from image API: %v", err)
	}

	if len(annotations) == 0 {
		return "", errors.New("image contained no text")
	}

	s := annotations[0].Description
	s = strings.Split(s, "\n")[0]

	debugPath := "debug/" + s
	debug, err := os.Create(debugPath)
	if err != nil {
		log.Fatalf("failed to open debug file: %v", err)
	}

	_, err = debug.WriteString(annotations[0].Description)
	if err != nil {
		log.Fatalf("failed to write to debug file: %v", err)
	}

	return s, nil
}

func ValidCard(cardName string) bool {
	baseUrl := "https://api.scryfall.com/cards/named"
	v := url.Values{}
	v.Set("exact", cardName)
	query := baseUrl + "?" + v.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
