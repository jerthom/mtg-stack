package identify

import (
	"context"
	"errors"
	"fmt"
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
	return s, nil
}
