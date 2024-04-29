package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

func CardTitle(imgName string) string {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	imgPath := "img/" + imgName
	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to get image from reader: %v", err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to get text from image API: %v", err)
	}

	if len(annotations) == 0 {
		fmt.Println("No text found")
	}

	s := annotations[0].Description
	s = strings.Split(s, "\n")[0]
	return s
}

func main() {
	entries, err := os.ReadDir("./img")
	if err != nil {
		log.Fatalf("Unable to read dir entries: %v", err)
	}

	for _, e := range entries {
		cardName := CardTitle(e.Name())
		fmt.Println(cardName)
	}
}
