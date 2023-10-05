package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
)

func downloadVideo(url string, fileName string, errCh chan error) {
	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		errCh <- err
		return
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		errCh <- err
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func clipVideo(inputFile string, outputFile string, start string, end string, errCh chan error) {
	cmd := exec.Command("ffmpeg", "-ss", start, "-i", inputFile, "-to", end, "-c", "copy", outputFile)
	err := cmd.Run()
	errCh <- err
}

func main() {
	url := flag.String("url", "", "YouTube video URL")
	timePairs := flag.String("timePairs", "", "Comma-separated time pairs in hh:mm:ss-hh:mm:ss format")
	flag.Parse()

	if *url == "" || *timePairs == "" {
		fmt.Println("Please provide a YouTube URL and time periods")
		return
	}

	pairs := strings.Split(*timePairs, ",")

	downloadFile := "temp_video.mp4"

	errCh := make(chan error, len(pairs)+1) // +1 for the download operation

	fmt.Println("Downloading video...")
	go downloadVideo(*url, downloadFile, errCh)

	var wg sync.WaitGroup

	for i, pair := range pairs {
		wg.Add(1)

		go func(i int, pair string) {
			defer wg.Done()

			times := strings.Split(pair, "-")
			if len(times) != 2 {
				fmt.Printf("Invalid time pair: %s\n", pair)
				return
			}

			start := times[0]
			end := times[1]
			clipFile := fmt.Sprintf("clipped_video_%d.mp4", i+1)

			fmt.Printf("Clipping video for period %s to %s...\n", start, end)
			clipVideo(downloadFile, clipFile, start, end, errCh)
		}(i, pair)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Printf("An error occurred: %v", err)
		}
	}

	fmt.Println("Done!")
}
